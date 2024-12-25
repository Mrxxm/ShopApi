package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net/http"
	"shop_api/user-web/forms"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano()) // 随机种子

	var code strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&code, "%d", numeric[rand.Intn(r)])
	}

	return code.String()
}

func SendSms(ctx *gin.Context) {

	// 1.表单验证
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	var mobile string = sendSmsForm.Mobile
	var code string = GenerateSmsCode(6)

	// 2.发送短信逻辑
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", "", "")
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = mobile                     // 手机号
	request.QueryParams["SignName"] = "啃肉提醒"                         // 阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_476685297"            // 阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + code + "}" // 短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
	// 3.将验证码保存起来
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	rdb.Set(context.Background(), mobile, code, 300*time.Second) // 300秒

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
