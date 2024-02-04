package middleware

import (
	"LanshanTeam-Examine/caller/model"
	"LanshanTeam-Examine/caller/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Mycliam struct {
	Username string `json:"username,omitempty"`
	jwt.StandardClaims
}

var secret = []byte("LanShanTeam examine secrete")

func GetToken(user *model.Userinfo) (string, error) {
	claims := Mycliam{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    "admin",
		},
	}
	badToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := badToken.SignedString(secret)
	return token, err
}

func ParseJWT(tokenstring string) (*Mycliam, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Mycliam{}, func(token *jwt.Token) (interface{}, error) {
		//解析方法使用此回调函数来提供用于验证的密钥。函数接收已解析但未验证的令牌。
		//这允许您在标识要使用哪个密钥的令牌（如“kid”）的标头。
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		if token.Valid && token.Claims.(*Mycliam) != nil {
			return token.Claims.(*Mycliam), nil
		}
	}
	return nil, err
}
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := ctx.Request.Header.Get("Authorization")
		utils.ClientLogger.Debug("TOKEN : " + value)
		tokenstr := strings.SplitN(value, " ", 2)
		if tokenstr[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "JWT 格式不正确",
			})
			ctx.Abort()
			return
		}
		if tokenstr[1] == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "JWT 为空",
			})
			ctx.Abort()
			return
		}
		cliam, err := ParseJWT(tokenstr[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "解析失败",
			})
			ctx.Abort()
			return
		} else if cliam.ExpiresAt < time.Now().Unix() {
			ctx.JSON(400, gin.H{
				"error": "token 超时",
			})
			ctx.Abort()
			return
		}

		ctx.Set("username", cliam.Username)
		ctx.Next()
	}
}
