package handle

import (
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
	user := &db.UserInfo{
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: int(req.PhoneNumber),
		Email:       req.Email,
	}
	//
	var realUser *db.UserInfo
	if err := checkUserIsAlreadyExist(user, realUser); err != nil {
		return &pb.RegisterResp{
			Message: "user already exist",
			Flag:    false,
		}, nil
	}
	if err := user.Create(); err != nil {
		utils.UserLogger.Error("ERROR:" + err.Error())
		return &pb.RegisterResp{
			Message: err.Error(),
			Flag:    false,
		}, err
	}
	utils.UserLogger.Info("user create success")
	return &pb.RegisterResp{
		Message: "user create success",
		Flag:    true,
	}, nil
}
func checkUserIsAlreadyExist(user *db.UserInfo, realUser *db.UserInfo) error {
	if err := user.Get("username", user.Username, realUser); err == nil {
		utils.UserLogger.Info("user already create,please change another user name")
		return errors.New("user already create,please change another user name")
	}
	return nil
}
