package handle

import (
	"LanshanTeam-Examine/server/user/dao/cathe"
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"errors"
	"log"
	"strconv"
)

type UserServer struct {
	pb.UnimplementedUserServer
}

func (U *UserServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	//先写一个密码登陆
	utils.UserLogger.Debug("username : " + req.GetUsername())
	utils.UserLogger.Debug("password : " + req.GetPassword())
	dbUser := &db.UserInfo{
		Username:    req.GetUsername(),
		Password:    req.GetPassword(),
		PhoneNumber: int(req.GetPhoneNumber()),
		Email:       req.GetEmail(),
	}
	catheUser := &cathe.UserInfoInCathe{
		*dbUser,
	}

	var realUser = new(db.UserInfo)
	//true nil 才是存在
	flag, err := catheCheckUserIsAlreadyExist(ctx, catheUser)
	if flag && err == nil {
		return &pb.RegisterResp{
			Message: "user already exist",
			Flag:    false,
		}, nil
	} else if !flag {
		if err := dbCheckUserIsAlreadyExist(ctx, catheUser, realUser); err == nil {
			return &pb.RegisterResp{
				Message: "user already exist",
				Flag:    false,
			}, nil
		}
	}
	//验证成功，用户不存在，准备创建用户，开始加盐加密
	dbUser.Password = utils.Encrypt(dbUser.Password)
	catheUser.Password = dbUser.Password //同步加盐加密
	//存储在数据库
	if err := dbUser.Create(); err != nil {
		utils.UserLogger.Error("ERROR:" + err.Error())
		return &pb.RegisterResp{
			Message: err.Error(),
			Flag:    false,
		}, err
	}
	utils.UserLogger.Info("user create success")
	//创建成功后，加载到缓存中
	err = catheUser.CreateUser(ctx)
	if err != nil {
		utils.UserLogger.Error("ERROR:" + err.Error())
		//创建用户成功了，只不过是redis无法创建缓存，所以不返回报错信息，只打日志
	} else {
		utils.UserLogger.Info("cathe: userinfo create success")

	}
	return &pb.RegisterResp{
		Message: "user create success",
		Flag:    true,
	}, nil
}

func (U *UserServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	utils.UserLogger.Debug("username : " + req.GetUsername())
	utils.UserLogger.Debug("password : " + req.GetPassword())
	//初始化环节
	dbUser := &db.UserInfo{
		Username:     req.GetUsername(),
		Password:     req.GetPassword(),
		PhoneNumber:  int(req.GetPhoneNumber()),
		Email:        req.GetEmail(),
		IsGithubUser: req.GetIsGithubUser(),
	}
	catheUser := &cathe.UserInfoInCathe{
		*dbUser,
	}
	var realUser = new(db.UserInfo)

	//github第三方登陆到逻辑
	if dbUser.IsGithubUser {
		err := GithubUserLogin(ctx, catheUser, realUser)
		if err != nil {
			utils.UserLogger.Error("github user login failed , ERROR:" + err.Error())
			return &pb.LoginResp{
				Message: "create github user information failed",
				Flag:    false,
			}, err
		}
		utils.UserLogger.Info("github user login success")
		return &pb.LoginResp{
			Message: "github user login success",
			Flag:    true,
		}, nil
	}

	//检查用户是否存在
	flag, err := catheCheckUserIsAlreadyExist(ctx, catheUser)
	if !flag && err != nil {
		if err := dbCheckUserIsAlreadyExist(ctx, catheUser, realUser); err != nil {
			return &pb.LoginResp{
				Message: "user not found",
				Flag:    false,
			}, nil
		}
	}
	log.Println("用户存在")

	//电话号码登陆
	if req.PhoneNumber != 0 {
		err := LoginByPhoneNumber(ctx, catheUser, realUser)
		if err != nil {
			utils.UserLogger.Error("login by phone number failed , ERROR:" + err.Error())
			return &pb.LoginResp{
				Message: "login by phone number failed",
				Flag:    false,
			}, err
		}
		utils.UserLogger.Info("github user login success")
		return &pb.LoginResp{
			Message: "login by phone number success",
			Flag:    true,
		}, nil
	}

	//用户存在 检查密码
	passwordInCathe, err := catheUser.GetWhat(ctx, "password")
	if err != nil {
		err = dbUser.Get("username", dbUser.Username, realUser)
		if err != nil {
			utils.UserLogger.Error("get password from database failed,ERROR:" + err.Error())
			return &pb.LoginResp{
				Message: "can't get user information from server",
				Flag:    false,
			}, err
		}

		err = utils.Compare(realUser.Password, catheUser.Password)

		if err != nil {
			utils.UserLogger.Error("wrong password in database,ERROR:" + err.Error())
			return &pb.LoginResp{
				Message: "wrong password",
				Flag:    false,
			}, errors.New("wrong password")
		}
		//存入缓存
		catheUser.UserInfo = *realUser
		err = catheUser.CreateUser(ctx)
		if err != nil {
			utils.UserLogger.Error("create user in cathe failed:" + err.Error())
		} else {
			utils.UserLogger.Info("create user in cathe success")
		}
	} else {
		//缓存命中
		err = utils.Compare(passwordInCathe, dbUser.Password)
		if err != nil {
			utils.UserLogger.Error("wrong password in cathe,ERROR:" + err.Error())
			return &pb.LoginResp{
				Message: "wrong password",
				Flag:    false,
			}, errors.New("wrong password")
		}
	}

	//登陆成功
	return &pb.LoginResp{
		Message: "login success",
		Flag:    true,
	}, nil
}

func GithubUserLogin(ctx context.Context, githubUser *cathe.UserInfoInCathe, realUser *db.UserInfo) error {
	//客户端传来登陆成功的用户名，我们先在缓存中找，没找到，到数据库找，没找到，就在数据库创建，再同步到缓存中
	//统一在用户名后+(github)标识
	githubUser.Username = githubUser.Username + "(github)"
	utils.UserLogger.Debug("github user name rename : " + githubUser.Username)
	flag, err := catheCheckUserIsAlreadyExist(ctx, githubUser)
	if flag && err == nil {
		utils.UserLogger.Debug("github user info found in cathe")
		return nil
	} else {
		err = dbCheckUserIsAlreadyExist(ctx, githubUser, realUser)
		if err == nil {
			utils.UserLogger.Debug("github user info fund in database")
			return nil
		}

	}
	//
	utils.UserLogger.Debug("github user info not fund in database,begin to create user")
	err = githubUser.Create()
	if err != nil {
		return err
	}
	githubUser.CreateUser(ctx) //不处理错误
	return nil
}

func LoginByPhoneNumber(ctx context.Context, enteredUser *cathe.UserInfoInCathe, realUser *db.UserInfo) error {
	//此时用户名是存在的，就只需要查询该电话号码跟数据库里面的电话号码是否相同。两种情况：1.用户没有绑定电话号码，为0；2.用户的电话号码绑定的不是这个
	//先查缓存，命中则返回，未命中就查数据库，查到就同步到缓存，并返回
	num, err := enteredUser.GetWhat(ctx, "phone_number")
	if err != nil {
		utils.UserLogger.Info("user register by phone number found in cathe ")
		err = enteredUser.Get("username", enteredUser.Username, realUser)
		if err != nil {
			utils.UserLogger.Error("login by phone number : database wrong")
			return err
		}
		if enteredUser.PhoneNumber != realUser.PhoneNumber {
			utils.UserLogger.Info("phone number is wrong compared with database")
			return errors.New("phone number is wrong compared with database")
		}
		err := enteredUser.CreateUser(ctx) //create info in cathe
		if err != nil {
			utils.UserLogger.Info("create user cathe failed , ERROR:" + err.Error())
		} else {
			utils.UserLogger.Info("create user cathe success ")
		}
		return nil

	} else {
		number, _ := strconv.Atoi(num)
		if number != enteredUser.PhoneNumber {
			utils.UserLogger.Error("phone number is wrong compared with cathe")
			return errors.New("phone number is wrong compared with cathe")
		}
		return nil
	}
}

func dbCheckUserIsAlreadyExist(ctx context.Context, user *cathe.UserInfoInCathe, realUser *db.UserInfo) error {

	if err := user.Get("username", user.Username, realUser); err != nil {
		utils.UserLogger.Info(" not fund in database")
		return errors.New("not fund in database")
	}
	user.Password = realUser.Password
	err := user.CreateUser(ctx) //create info in cathe
	if err != nil {
		utils.UserLogger.Info("create user cathe failed , ERROR:" + err.Error())
	} else {
		utils.UserLogger.Info("create user cathe success ")
	}
	return nil
}
func catheCheckUserIsAlreadyExist(ctx context.Context, catheUser *cathe.UserInfoInCathe) (bool, error) {
	_, err := catheUser.GetWhat(ctx, "username") //先找，没找到就是没存到缓存,找到就是存在
	if err != nil {
		utils.UserLogger.Info("not found in cathe ERROR:" + err.Error())
		return false, errors.New("not found in cathe")
	}
	utils.UserLogger.Info("found in cathe")
	return true, nil
}
