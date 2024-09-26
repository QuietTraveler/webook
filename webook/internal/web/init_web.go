package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
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

	registerUsersRoutes(server, u)
	return
}
func registerUsersRoutes(server *gin.Engine, u *UserHandler) {
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.GET("/users/profile", u.Profile)
}
