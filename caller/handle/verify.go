package handle

import (
	"LanshanTeam-Examine/caller/model"
	"LanshanTeam-Examine/caller/pkg/consts"
	"LanshanTeam-Examine/caller/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"strings"
	"time"
)

// 一个电话号码对应它的验证码信息
var CodeInfo = make(map[string]model.CodeInfo)
var aliyun model.Config

func SendCode(c *gin.Context) {
	viper.SetConfigFile("./config/gameConfig.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		utils.ClientLogger.Error("sent code can't read in config file , ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "server unavailable,服务不可用",
			"error":   err.Error(),
		})
		return
	}
	err = viper.Unmarshal(&aliyun)
	if err != nil {
		utils.ClientLogger.Error("sent code can't unmarshal the config , ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "server unavailable,服务不可用",
			"error":   err.Error(),
		})
		return
	}

	client, err := CreateClient(tea.String(aliyun.AccessKeyId), tea.String(aliyun.AccessKeySecret))

	if err != nil {
		utils.ClientLogger.Error("sent code can't create caller , ERROR:" + err.Error())
		c.JSON(400, gin.H{
			"code":    consts.ServeUnavailable,
			"message": "server unavailable,服务不可用",
			"error":   err.Error(),
		})
		return
	}

	phoneNumber := c.PostForm("phone_number")

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("lance47"),
		TemplateCode:  tea.String("SMS_465025228"),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(GenerateCode(phoneNumber)),
	}
	utils.ClientLogger.Debug("the params :::" + GenerateCode(phoneNumber))
	runtime := &util.RuntimeOptions{}
	tryErr := func() (err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}
		utils.ClientLogger.Info("INFO,message : " + *resp.Body.Message)
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		str, _ := util.AssertAsString(error.Message)
		utils.ClientLogger.Debug("the error of failed to send code :" + tea.StringValue(str))
		c.JSON(400, gin.H{
			"code":    consts.SendCodeFailed,
			"message": "send code failed",
			"error":   str,
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    consts.SendCodeSuccess,
		"message": "sent code success",
		"error":   nil,
	})
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// 生成验证码
func GenerateCode(num string) string {
	number := rand.New(rand.NewSource(time.Now().Unix()))
	var newCodeInfo model.CodeInfo
	newCodeInfo.Code = number.Int63n(9000) + 1000
	newCodeInfo.ExpireAt = time.Now().Add(10 * time.Minute)
	CodeInfo[num] = newCodeInfo
	return fmt.Sprintf("{\"code\":\"%d\"}", CodeInfo[num].Code)
}

// 校验验证码
func VerifyCode(code int64, num string) error {
	val, ok := CodeInfo[num]
	log.Println("验证码：", val.Code)
	if !ok {
		utils.ClientLogger.Info("not found code ,should generate code ")
		return errors.New("not found code")
	}
	if time.Now().After(val.ExpireAt) {
		utils.ClientLogger.Info("the code is expired")
		return errors.New("code is expired")
	}
	if val.Code != code {
		utils.ClientLogger.Info("this code is wrong")
		return errors.New("code is wrong")
	}
	return nil
}
