package handle

import (
	"LanshanTeam-Examine/server/user/dao/cathe"
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"errors"
	"log"
)

type UserServer struct {
	pb.UnimplementedUserServer
}

func (U *UserServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	//先写一个密码登陆
	utils.UserLogger.Debug("username : " + req.Username)
	utils.UserLogger.Debug("password : " + req.Password)
	dbUser := &db.UserInfo{
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: int(req.PhoneNumber),
		Email:       req.Email,
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

func (u *UserServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	utils.UserLogger.Debug("username : " + req.Username)
	utils.UserLogger.Debug("password : " + req.Password)

	dbUser := &db.UserInfo{
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: int(req.PhoneNumber),
		Email:       req.Email,
	}
	catheUser := &cathe.UserInfoInCathe{
		*dbUser,
	}

	var realUser = new(db.UserInfo)
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
		log.Println("到此一游")
		log.Println(realUser.Password)
		log.Println(catheUser.Password)
		log.Println(dbUser.Password)
		log.Println("===========================")
		err = utils.Compare(realUser.Password, catheUser.Password)

		if err != nil {
			utils.UserLogger.Error("wrong password,ERROR:" + err.Error())
			return &pb.LoginResp{
				Message: "wrong password",
				Flag:    false,
			}, errors.New("wrong password")
		}
		//存入缓存
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
			utils.UserLogger.Error("wrong password,ERROR:" + err.Error())
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

func dbCheckUserIsAlreadyExist(ctx context.Context, user *cathe.UserInfoInCathe, realUser *db.UserInfo) error {

	if err := user.Get("username", user.Username, realUser); err != nil {
		utils.UserLogger.Info(" not fund in database")
		return errors.New("not fund in database")
	}
	err := user.CreateUser(ctx)
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
		utils.UserLogger.Info("not found in cathe ")
		return false, errors.New("not found in cathe")
	}
	utils.UserLogger.Info("found in cathe")
	return true, nil
}
