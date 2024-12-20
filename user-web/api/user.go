package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/user-web/global"
	"shop_api/user-web/global/response"
	"shop_api/user-web/proto"
	"time"
)

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		global.GetSugar().Errorw("[HandleGrpcErrorToHttp] 处理grpc错误", "msg", err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			case codes.Unavailable:
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "服务不可用"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			}
		}
		return
	}
}

func GetUserList(ctx *gin.Context) {
	global.Sugar.Debug("获取用户列表页")

	// 1.拨号连接grpc服务
	connect, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		global.GetSugar().Errorw("[GetUserList] 连接 [user_srv] 失败", "msg", err.Error())
	}
	defer connect.Close()

	// 2.生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(connect)
	UserListResponse, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{Page: 1, PageSize: 5})
	if err != nil {
		global.GetSugar().Errorw("[GetUserList] 调用 [user_srv] GetUserList 失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 3.返回数据
	//result := make(map[string]interface{}, 0)
	result := make([]interface{}, 0)
	for _, value := range UserListResponse.Data {

		fmt.Println(value.Birthday)

		userInfo := response.UserResponse{
			Id:       value.Id,
			Nickname: value.Nickname,
			//Birthday: time.Time(time.Unix(int64(value.Birthday), 0)).Format("2006-01-02 15:04:05"),
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Mobile:   value.Mobile,
			Gender:   value.Gender,
		}

		//userInfo := map[string]interface{}{
		//	"id":       value.Id,
		//	"name":     value.Nickname,
		//	"birthday": value.Birthday,
		//	"gender":   value.Gender,
		//	"mobile":   value.Mobile,
		//}
		result = append(result, userInfo)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     result,
		"total":    UserListResponse.Total,
		"page":     1,
		"pageSize": 5,
	})
}
