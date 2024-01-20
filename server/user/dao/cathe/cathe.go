package cathe

import (
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Test() {
	name := RedisClient.Get(context.Background(), "name")
	utils.UserLogger.Info(name.Val())
	utils.UserLogger.Debug(name.String())
}
