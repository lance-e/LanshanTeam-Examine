package userModule

import (
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	UserClient pb.UserClient

	UserConn *grpc.ClientConn
)

func UserRPC() {

	var opts []grpc.DialOption
	var err error
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	UserConn, err = grpc.Dial("localhost:10001", opts...)

	if err != nil {
		utils.ClientLogger.Debug(err.Error())
		fmt.Println("TEST CONNECTION...........ERROR")
	}
	UserClient = pb.NewUserClient(UserConn)

}
