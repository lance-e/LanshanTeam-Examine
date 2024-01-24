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
	result, err := RedisClient.HGet(ctx, "\"userinfo:"+b.Username+"\"", "\""+what+"\"").Result()
	if err != nil {
		utils.UserLogger.Error("redis error HGET:" + err.Error())
		return "", err
	}
	return result, nil
}
func (b *UserInfoInCathe) CreateUser(ctx context.Context) error {
	err := RedisClient.HSet(ctx, "userinfo:"+b.Username, "username", b.Username, "password", b.Password, "phone_number", b.PhoneNumber, "email", b.Email, "is_github_user", b.IsGithubUser).Err()
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
