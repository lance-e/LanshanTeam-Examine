package handle

import (
	"LanshanTeam-Examine/client/model"
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

var Info struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
}

type username struct {
	Username string `json:"login"`
}

var config model.Config

func GithubRegisterAndLogin(c *gin.Context) {

	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		utils.ClientLogger.Error("GithubRegisterAndLogin can't read in config file , ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "server unavailable,服务不可用",
			"error":   err.Error(),
		})
		return
	}
	log.Println(config)
	err = viper.Unmarshal(&config)
	log.Println(config)

	if err != nil {
		utils.ClientLogger.Error("GithubRegisterAndLogin can't unmarshal the config , ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "server unavailable,服务不可用",
			"error":   err.Error(),
		})
		return
	}

	//开始调api
	callback := "http://localhost:8080/user/githubCallback"
	c.Redirect(307, "https://github.com/login/oauth/authorize?client_id="+config.ClientId+"&redirect_url="+callback+"&scope=user&state=random")
}
func GithubCallback(c *gin.Context) {

	var info = Info
	code := c.Query("code")
	utils.ClientLogger.Info("github code in query :" + code)
	state := c.Query("state")
	utils.ClientLogger.Info("github state in query : " + state)
	//url to get token
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		config.ClientId, config.ClientSecrets, code,
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		utils.ClientLogger.Error("REQUEST: get github token error:" + err.Error())
		return
	}
	req.Header.Set("Accept", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ClientLogger.Error("get token request failed, ERROR:" + err.Error())
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(data, &info)

	if err != nil {
		utils.ClientLogger.Error("can't unmarshal" + err.Error())
		return
	}

	client = *http.DefaultClient
	req, err = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		utils.ClientLogger.Error("REQUEST:exchange userinfo from token error ,ERROR : " + err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer "+info.Token)

	resp, err = client.Do(req)
	if err != nil {
		utils.ClientLogger.Error("exchange userinfo request failed, " + err.Error())
		return
	}
	defer resp.Body.Close()

	data, _ = io.ReadAll(resp.Body)
	var name username
	json.Unmarshal(data, &name)
	utils.ClientLogger.Info("username use github  : " + name.Username)

	//开始调用远程注册服务
	//userModule.UserClient.Register(c,)

	c.Redirect(307, "http://localhost:8080/") //回调到主页
}
