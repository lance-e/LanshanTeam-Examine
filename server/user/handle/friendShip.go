package handle

import (
	"LanshanTeam-Examine/server/user/dao/cathe"
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"errors"
)

func (U *UserServer) AddFriend(ctx context.Context, req *pb.AddFriendReq) (*pb.AddFriendResp, error) {
	//先检查有没有向目标用户发送过好友申请，没有就发送，有就返回
	friendship := &db.FriendShip{
		Sender:   req.Sender,
		Receiver: req.Receiver,
	}
	if req.IsRequest {
		//检查接受者是否存在
		catheUser := cathe.UserInfoInCathe{}
		catheUser.Username = req.Receiver
		flag, err := catheCheckUserIsAlreadyExist(ctx, &catheUser)
		if !flag && err != nil {
			if err := dbCheckUserIsAlreadyExist(ctx, &catheUser, nil); err != nil {
				return &pb.AddFriendResp{
					Flag: false,
				}, errors.New("receiver not found ")
			}
		}

		//查数据库，懒得弄缓存了
		flag = friendship.IsRequestAlreadyExists()
		if flag {
			return &pb.AddFriendResp{
				Flag: false,
			}, errors.New("request already exists or database was error")
		}

		err = friendship.Create()
		if err != nil {
			utils.UserLogger.Error("create friend ship failed , ERROR:" + err.Error())
			return &pb.AddFriendResp{
				Flag: false,
			}, errors.New("database  error")
		}

		return &pb.AddFriendResp{
			Flag: true,
		}, nil

	} else {
		err := friendship.Update()
		if err != nil {
			return &pb.AddFriendResp{
				Flag: false,
			}, err
		}
		return &pb.AddFriendResp{
			Flag: true,
		}, nil
	}

}
