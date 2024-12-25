package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/user-web/forms"
	"shop_api/user-web/global"
	"shop_api/user-web/global/response"
	"shop_api/user-web/middlewares"
	"shop_api/user-web/models"
	"shop_api/user-web/proto"
	"strconv"
	"time"
)

func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"msg": errs.Translate(global.Trans),
	})
}

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
	global.Sugar.Debug("获取用户列表页接口")

	claims, _ := ctx.Get("claims")
	userId := claims.(*models.CustomClaims).ID
	global.Sugar.Debug("用户id:", userId)

	// 1.拨号连接grpc服务
	connect, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		global.GetSugar().Errorw("[GetUserList] 连接 [user_srv] 失败", "msg", err.Error())
	}
	defer connect.Close()

	// 2.生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(connect)
	page := ctx.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := ctx.DefaultQuery("page_size", "5")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	UserListResponse, err := userSrvClient.GetUserList(
		context.Background(),
		&proto.PageInfo{
			Page:     uint32(pageInt),
			PageSize: uint32(pageSizeInt),
		})
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

func PasswordLogin(ctx *gin.Context) {
	global.Sugar.Debug("登录接口")
	// 1.表单验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 2.验证码验证
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "验证码错误"})
		return
	}

	// 3.拨号连接grpc服务
	connect, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		global.GetSugar().Errorw("[GetUserList] 连接 [user_srv] 失败", "msg", err.Error())
	}
	defer connect.Close()

	// 4.生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(connect)
	UserInfoResponse, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 5.验证密码是否正确
	CheckResponse, err := userSrvClient.CheckPassword(ctx, &proto.PasswordCheckInfo{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: UserInfoResponse.Password,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if !CheckResponse.Success {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "密码错误"})
		return
	}

	// 6.返回数据
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(UserInfoResponse.Id),
		NickName:    UserInfoResponse.Nickname,
		AuthorityId: uint(UserInfoResponse.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),            // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24, // 过期时间 一天
			Issuer:    "xxm",                        //签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "内部错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

func Register(ctx *gin.Context) {
	// 1.表单验证
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 2.验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	redisValue := rdb.Get(context.Background(), registerForm.Mobile).Val()
	if redisValue != registerForm.Code {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "验证码错误"})
		return
	}

	// 3.用户注册
	connect, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		global.GetSugar().Errorw("[Register] 连接 [user_srv] 失败", "msg", err.Error())
	}
	defer connect.Close()

	// 4.生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(connect)
	UserInfoResponse, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   registerForm.Mobile,
		Nickname: "",
		Password: registerForm.Password,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "注册成功",
		"data": gin.H{
			"user_id": UserInfoResponse.Id,
		},
	})
}
