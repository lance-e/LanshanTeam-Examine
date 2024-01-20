package Init

import (
	"LanshanTeam-Examine/server/user/dao/db"
	"LanshanTeam-Examine/server/user/pkg/utils"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Info struct {
	DbInfo `mapstructure:"mysql" json:"mysql"`
	//catheInfo `mapstruct:"redis" json:"cathe_info"`
}

// 结构体名要大写
type DbInfo struct {
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Dbname   string `mapstructure:"dbname" json:"dbname"`
}

//type catheInfo struct {
//}

var info Info
var err error

func InitMysql() {
	viper.SetConfigFile("./config/userConfig.yaml")
	err = viper.ReadInConfig()

	if err != nil {
		utils.UserLogger.Panic("read the mysql config file failed")

	}
	fmt.Println(viper.Get("mysql.username"))

	err = viper.Unmarshal(&info)

	if err != nil {
		utils.UserLogger.Panic(err.Error())
		return
	}

	dsn := strings.Join([]string{info.DbInfo.Username, ":", info.DbInfo.Password, "@tcp(", info.DbInfo.Host,
		":", info.DbInfo.Port, ")/", info.DbInfo.Dbname, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrate()
	if err != nil {
		utils.UserLogger.Panic("migrate failed" + err.Error())
	}
	fmt.Println("success")
}
