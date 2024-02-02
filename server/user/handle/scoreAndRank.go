package handle

import (
	"LanshanTeam-Examine/server/user/dao/cathe"
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pb"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
)

func (U *UserServer) AddScore(ctx context.Context, req *pb.AddScoreReq) (*pb.AddScoreResp, error) {
	name := req.Username
	if req.IsGithubName {
		name += "(github)"
	}
	user := db.UserInfo{
		Username: name,
	}
	err := user.AddScore()
	if err != nil {
		utils.UserLogger.Debug("addScore handle can't add score")
		return &pb.AddScoreResp{
			Message: "can't add score",
		}, err
	}
	//更新排行榜缓存
	catheUser := cathe.UserInfoInCathe{UserInfo: user}
	err = catheUser.UpdateRank(ctx)
	if err != nil {
		utils.UserLogger.Debug("update rank cathe failed")
	}
	utils.UserLogger.Debug("addScore handle  add score success")
	return &pb.AddScoreResp{
		Message: "add score success",
	}, nil
}
func (U *UserServer) Rank(ctx context.Context, req *pb.RankReq) (*pb.RankResp, error) {
	var rank []*pb.Rank
	result, err := cathe.GetRank(ctx)
	if err != nil || len(result) == 0 {
		utils.UserLogger.Debug("cathe had not anything about rank ")
		if err := MigrateIntoCathe(ctx); err != nil {
			utils.UserLogger.Error("can't migrate all information into redis")
			return &pb.RankResp{
				Rank:    nil,
				Message: "can't migrate all information into redis",
			}, err
		}
		result, err = cathe.GetRank(ctx)
		if err != nil {
			utils.UserLogger.Debug("what fxxk?")
			return &pb.RankResp{
				Rank:    nil,
				Message: "can't migrate all information into redis",
			}, err
		}
	}
	for _, v := range result {
		rank = append(rank, &pb.Rank{
			Username: v.Member.(string),
			Score:    int64(v.Score),
		})
	}
	return &pb.RankResp{
		Rank:    rank,
		Message: "the rank of score",
	}, nil
}

func MigrateIntoCathe(ctx context.Context) error {
	var infos []cathe.UserInfoInCathe
	err := db.DB.Model(&db.UserInfo{}).Select("username", "score").Find(&infos).Error
	if err != nil {
		utils.UserLogger.Error("can't get all the information")
		return err
	}
	for _, v := range infos {
		err = v.UpdateRank(ctx)
		if err != nil {
			utils.UserLogger.Error("can't migrate this information")
		}
	}
	return nil
}
