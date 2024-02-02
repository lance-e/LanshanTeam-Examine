package cathe

import (
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisClient *redis.Client

type UserInfoInCathe struct {
	db.UserInfo
}

// GetWhat 缓存操作，
func (b *UserInfoInCathe) GetWhat(ctx context.Context, what string) (string, error) {
	result, err := RedisClient.HGet(ctx, "userinfo:"+b.Username+"", ""+what+"").Result()
	if err != nil {
		utils.UserLogger.Error("redis error HGET:" + err.Error())
		return "", err
	}
	return result, nil
}
func (b *UserInfoInCathe) GetAll(ctx context.Context) (map[string]string, error) {
	result, err := RedisClient.HGetAll(ctx, "userinfo:"+b.Username+"").Result()
	if err != nil {
		utils.UserLogger.Error("redis error HGETALL:" + err.Error())
		return nil, err
	}
	return result, err
}

func (b *UserInfoInCathe) CreateUser(ctx context.Context) error {
	err := RedisClient.HSet(ctx, "userinfo:"+b.Username, "username", b.Username,
		"password", b.Password, "phone_number", b.PhoneNumber, "email", b.Email, "is_github_user", b.IsGithubUser, "score", b.Score).Err()
	SetExpireTime(ctx, "userinfo:"+b.Username)
	if err != nil {
		utils.UserLogger.Error("redis error HSET:" + err.Error())
		return err
	}
	return nil
}

func SetExpireTime(ctx context.Context, key string) {
	err := RedisClient.Expire(ctx, key, 10*time.Minute).Err()
	if err != nil {
		utils.UserLogger.Error("set expire time EXPIRE: " + err.Error())
		return
	}
}
func (b *UserInfoInCathe) UpdateRank(ctx context.Context) error {
	err := RedisClient.ZAdd(ctx, "rank", redis.Z{Score: float64(b.Score), Member: b.Username}).Err()
	if err != nil {
		utils.UserLogger.Error("redis error ZADD:" + err.Error())
		return err
	}
	utils.UserLogger.Debug("update rank success")
	return nil
}
func GetRank(ctx context.Context) ([]redis.Z, error) {
	result, err := RedisClient.ZRevRangeWithScores(ctx, "rank", 0, 9).Result()
	if err != nil {
		utils.UserLogger.Error("redis error Zrevrange:" + err.Error())
		return nil, err
	}
	utils.UserLogger.Debug("Zrevrange success")
	return result, nil
}
