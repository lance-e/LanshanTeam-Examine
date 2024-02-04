package userModule

import (
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/discovery"
	"LanshanTeam-Examine/caller/rpc/userModule/pb"
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
	var dis *discovery.EtcdDiscovery

	//服务注册
	dis, err = discovery.NewServerDiscovery([]string{"localhost:2379"})
	if err != nil {
		utils.ClientLogger.Debug("can't init a new discovery")
		return
	}
	defer dis.Close()
	err = dis.ServerDiscovery("user")
	if err != nil {
		utils.ClientLogger.Debug("can't discover the service")
		return
	}

	addr, err := dis.GetService("user")

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	UserConn, err = grpc.Dial(addr, opts...)

	if err != nil {
		utils.ClientLogger.Debug(err.Error())
		fmt.Println("TEST CONNECTION...........ERROR")
	}
	UserClient = pb.NewUserClient(UserConn)

}
