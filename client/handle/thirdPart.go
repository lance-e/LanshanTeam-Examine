package handle

import (
	"LanshanTeam-Examine/client/api/middleware"
	"LanshanTeam-Examine/client/model"
	"LanshanTeam-Examine/client/pkg/consts"
	"LanshanTeam-Examine/client/pkg/utils"
	"LanshanTeam-Examine/client/rpc/userModule"
	"LanshanTeam-Examine/client/rpc/userModule/pb"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
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

var github model.Config

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
	err = viper.Unmarshal(&github)

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
	c.Redirect(307, "https://github.com/login/oauth/authorize?client_id="+github.ClientId+"&redirect_url="+callback+"&scope=user&state=random")
}
func GithubCallback(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(400, gin.H{
				"code":    consts.ServeUnavailable,
				"message": "github login serve unavailable",
				"error":   err.(string),
			})
			return
		}
	}()

	var info = Info
	code := c.Query("code")
	utils.ClientLogger.Info("github code in query :" + code)
	state := c.Query("state")
	utils.ClientLogger.Info("github state in query : " + state)
	//url to get token
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		github.ClientId, github.ClientSecrets, code,
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
	githubUserName := name.Username
	//开始调用远程登陆服务
	loginreq := &pb.LoginReq{
		Username:     githubUserName,
		Password:     "",
		IsGithubUser: true,
	}
	loginResp, err := userModule.UserClient.Login(c, loginreq)

	utils.ClientLogger.Debug("request send")
	if err != nil {
		utils.ClientLogger.Error("create github user information failed")
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": loginResp.GetMessage(),
			"error":   err.Error(),
		})
		//return //不返回,因为想做重定向，返回了，就无法重定向了
	} else {
		utils.ClientLogger.Error("create github user information success")
		var user model.Userinfo
		user.Username = githubUserName
		token, err := middleware.GetToken(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"code":    consts.GenerateTokenFailed,
				"message": "generate token failed",
				"error":   err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"code":    consts.LoginSuccess,
			"message": "login success",
			"error":   "",
			"token":   token,
		})
	}
	//c.Redirect(307, "http://localhost:8080/") //回调到主页

}
