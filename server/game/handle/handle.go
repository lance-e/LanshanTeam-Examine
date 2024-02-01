package handle

import (
	"LanshanTeam-Examine/server/game/dao/db"
	"LanshanTeam-Examine/server/game/pb"
	"LanshanTeam-Examine/server/game/utils"
	"context"
)

type GameServer struct {
	pb.UnimplementedGameServer
}

func (g *GameServer) Save(ctx context.Context, req *pb.SaveReq) (*pb.SaveResp, error) {
	var step = &db.GameSteps{
		RoomHost: req.RoomHost,
		Player:   req.Player,
		Row:      req.Row,
		Column:   req.Column,
	}

	err := step.Create()
	if err != nil {
		utils.GameLogger.Debug("can't save the step")
		return &pb.SaveResp{
			Message: "can't save the step",
		}, err
	}

	return &pb.SaveResp{
		Message: "save success",
	}, nil
}
func (g *GameServer) ShowSteps(ctx context.Context, req *pb.ShowStepsReq) (*pb.ShowStepsReps, error) {
	var room = &db.GameSteps{
		RoomHost: req.RoomHost,
	}
	info, err := room.Get()
	if err != nil {
		utils.GameLogger.Debug("show steps failed ,error :" + err.Error())
		return &pb.ShowStepsReps{
			AllStep: make([]*pb.Step, 0),
		}, err
	}
	var resp []*pb.Step
	for _, v := range info {
		resp = append(resp, &pb.Step{
			RoomHost: v.RoomHost,
			Player:   v.Player,
			Row:      v.Row,
			Column:   v.Column,
		})
	}
	utils.GameLogger.Debug("show steps success")
	return &pb.ShowStepsReps{
		AllStep: resp,
	}, err
}
