package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"webook/webook/internal/web/middleware"
)

func RegisterRoutes(u *UserHandler) (server *gin.Engine) {
	server = gin.Default()

	//配置跨域请求
	//天坑，这里的Use方法需要放在注册路由之前，否则跨域请求会失败
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"}, // 使用明确的域名
		//AllowMethods:     []string{"POST", "GET"},    // 明确允许的 HTTP 方法
		AllowHeaders: []string{"Content-Type", "Authorization"}, // 确保必要的头部被允许
		//	ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true, // 允许 cookie
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//配置session中间件
	//用于解决HTTP协议的无状态性问题
	//
	//初始化一个cookie验证的会话存储器
	//由于会话数据的保密性需求，使用cookie.NewStore创建一个会话存储实例
	//使用"secret"作为加密密钥
	store := cookie.NewStore([]byte("secret"))
	// 配置会话中间件
	// sessions.Sessions方法创建一个名为"mysession"的会话
	// 使用上一行创建的store来存储会话数据
	// 通过server.Use将该会话中间件挂载到http服务器中
	// 以便每个请求都能进行会话管理
	server.Use(sessions.Sessions("webook", store))

	// 此部分用于验证用户是否已登录
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup", "/users/login").Build())

	registerUsersRoutes(server, u)
	return
}
func registerUsersRoutes(server *gin.Engine, u *UserHandler) {
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.GET("/users/profile", u.Profile)
}
