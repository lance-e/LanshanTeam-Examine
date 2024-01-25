package handle

import (
	"LanshanTeam-Examine/server/user/dao/cathe"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"strconv"
)

func (U *UserServer) HomePage(ctx context.Context, req *pb.HomePageReq) (*pb.HomePageResp, error) {
	var realUser cathe.UserInfoInCathe
	realUser.Username = req.Username
	utils.UserLogger.Debug("the user name HomePage function received:" + realUser.Username)
	info, err := realUser.GetAll(ctx)
	if err != nil {
		utils.UserLogger.Debug("not found in cathe ")
		err = realUser.Get("username", realUser.Username, &realUser.UserInfo)
		if err != nil {
			utils.UserLogger.Error("user information not found in database")
			return &pb.HomePageResp{
				Username: req.Username,
			}, err
		} else {
			utils.UserLogger.Error("user information  found in database")
			return &pb.HomePageResp{
				Username:    realUser.Username,
				PhoneNumber: int64(realUser.PhoneNumber),
				Email:       realUser.Email,
				Score:       int64(realUser.Score),
			}, nil
		}
	} else {
		num, _ := strconv.Atoi(info["phone_number"])
		score, _ := strconv.Atoi(info["score"])
		return &pb.HomePageResp{
			Username:    info["username"],
			PhoneNumber: int64(num),
			Email:       info["email"],
			Score:       int64(score),
		}, nil
	}

}
