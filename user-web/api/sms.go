package api

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
)

func SendSms(ctx *gin.Context) {

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
	request.QueryParams["PhoneNumbers"] = ""                             // 手机号
	request.QueryParams["SignName"] = "啃肉提醒"                             // 阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_476685297"                // 阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + "777777" + "}" // 短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
}
