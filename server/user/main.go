package main

import (
	"LanshanTeam-Examine/server/user/dao/Init"
	"LanshanTeam-Examine/server/user/dao/cathe"
)

func main() {
	Init.InitMysql()
	Init.InitRedis()
	cathe.Test()
}
