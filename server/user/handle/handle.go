package handle

import (
	"LanshanTeam-Examine/server/user/dao/cathe"
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"errors"
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

	var realUser *db.UserInfo

	flag, err := catheCheckUserIsAlreadyExist(ctx, catheUser, realUser)
	if err != nil {
		return &pb.RegisterResp{
			Message: "user already exist",
			Flag:    false,
		}, nil
	} else if !flag {
		if err := dbCheckUserIsAlreadyExist(dbUser, realUser); err != nil {
			return &pb.RegisterResp{
				Message: "user already exist",
				Flag:    false,
			}, nil
		}
	}
	//验证成功，用户不存在，准备创建用户，开始加盐加密
	dbUser.Password, err = utils.Encrypt(dbUser.Password)
	if err != nil {
		utils.UserLogger.Error("ERROR:" + err.Error())
		return &pb.RegisterResp{
			Message: err.Error(),
			Flag:    false,
		}, err
	}
	utils.UserLogger.Debug("HHHHHHHHHHHHH:cathe user' passwrod :" + catheUser.Password)
	catheUser.Password = dbUser.Password
	utils.UserLogger.Debug("XXXXXXXXXXXXX:cathe user' passwrod :" + catheUser.Password)

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
	}
	utils.UserLogger.Info("cathe: userinfo create success")

	return &pb.RegisterResp{
		Message: "user create success",
		Flag:    true,
	}, nil

}

func dbCheckUserIsAlreadyExist(user *db.UserInfo, realUser *db.UserInfo) error {

	if err := user.Get("username", user.Username, realUser); err == nil {
		utils.UserLogger.Info("user already create,please change another user name")
		return errors.New("user already create,please change another user name")
	}
	return nil
}
func catheCheckUserIsAlreadyExist(ctx context.Context, catheUser *cathe.UserInfoInCathe, realUser *db.UserInfo) (bool, error) {
	realUsername, err := catheUser.GetWhat(ctx, "username") //先找，没找到就是没存到缓存，找到之后再比较
	if err != nil {
		utils.UserLogger.Info("no found in cathe ")
		return false, nil
	}
	if realUsername != catheUser.Username {
		utils.UserLogger.Info("user already create,please change another user name")
		return false, errors.New("user already create,please change another user name")
	}
	return true, nil
}
