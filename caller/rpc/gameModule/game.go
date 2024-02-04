package gameModule

import (
	"LanshanTeam-Examine/caller/pkg/utils"
	"LanshanTeam-Examine/caller/rpc/discovery"
	"LanshanTeam-Examine/caller/rpc/gameModule/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GameClient pb.GameClient

	GameConn *grpc.ClientConn
)

func GameRPC() {
	var dis *discovery.EtcdDiscovery
	var opts []grpc.DialOption
	var err error

	dis, err = discovery.NewServerDiscovery([]string{"localhost:2379"})
	if err != nil {
		utils.ClientLogger.Debug("can't init a new discovery")
		return
	}
	defer dis.Close()
	err = dis.ServerDiscovery("game")
	if err != nil {
		utils.ClientLogger.Debug("can't discovery game service")
		return
	}

	addr, err := dis.GetService("game")
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	GameConn, err = grpc.Dial(addr, opts...)

	if err != nil {
		utils.ClientLogger.Debug(err.Error())
		fmt.Println("TEST CONNECTION...........ERROR")
	}
	GameClient = pb.NewGameClient(GameConn)

}
