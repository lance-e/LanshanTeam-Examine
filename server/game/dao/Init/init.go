package Init

import (
	"LanshanTeam-Examine/server/game/dao/cathe"
	"LanshanTeam-Examine/server/game/dao/db"
	"LanshanTeam-Examine/server/game/utils"

	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type Info struct {
	DbInfo    `mapstructure:"mysql" json:"mysql"`
	CatheInfo `mapstructure:"redis" json:"redis"`
}

// 结构体名要大写
type DbInfo struct {
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Dbname   string `mapstructure:"dbname" json:"dbname"`
}

type CatheInfo struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	Db       int    `mapstructure:"db" json:"db"`
}

var info Info
var err error

func readConfig() {
	viper.SetConfigFile("./config/gameConfig.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		utils.GameLogger.Panic("read the mysql config file failed")
	}
	err = viper.Unmarshal(&info)
	if err != nil {
		utils.GameLogger.Panic(err.Error())
	}
}

func InitMysql() {
	readConfig()
	dsn := strings.Join([]string{info.DbInfo.Username, ":", info.DbInfo.Password, "@tcp(", info.DbInfo.Host,
		":", info.DbInfo.Port, ")/", info.DbInfo.Dbname, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.GameLogger.Panic("couldn't open mysql")
	}
	err = db.Migrate()
	if err != nil {
		utils.GameLogger.Panic("migrate failed" + err.Error())
	}
	utils.GameLogger.Info("launch mysql successful")
}
func InitRedis() {
	readConfig()
	cathe.RedisClient = redis.NewClient(&redis.Options{
		Addr:     info.CatheInfo.Host + ":" + info.CatheInfo.Port,
		Password: info.CatheInfo.Password,
		DB:       info.CatheInfo.Db,
	})
	resp, err := cathe.RedisClient.Ping(context.Background()).Result()
	if err != nil {
		utils.GameLogger.Panic(err.Error())
	}
	if resp != "PONG" {
		utils.GameLogger.Error("ERROR:" + err.Error())
		return
	}
	utils.GameLogger.Info("launch redis success")
}
