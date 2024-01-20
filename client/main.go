package main

import "LanshanTeam-Examine/client/api/router"

func main() {
	engine := router.NewRouter()
	engine.Run("127.0.0.1:8080")
}
