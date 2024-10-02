package main

import (
	"github.com/ecodeclub/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
	"webook/webook/internal/web/middleware"
)

func initDB() (db *gorm.DB, err error) {
	dsn := "root:root@tcp(124.71.99.27:13316)/webook"
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		// 只在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化出错，应用就不要启动了
		return nil, err
	}

	if err := dao.InitTable(db); err != nil {
		return nil, err
	}
	return
}

func initUser(db *gorm.DB) (u *web.UserHandler) {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u = web.NewUserHandler(svc)
	return u
}

func initWebServer() (server *gin.Engine) {
	server = gin.Default()

	//配置跨域请求
	//天坑，这里的Use方法需要放在注册路由之前，否则跨域请求会失败
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"}, // 使用明确的域名
		//AllowMethods:     []string{"POST", "GET"},    // 明确允许的 HTTP 方法
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // 确保必要的头部被允许
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true, // 允许 cookie
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//限流
	redisClient := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "124.71.99.27:6379",
	})
	//创建限流器，设置时间窗口大小和请求阈值
	limiter := ratelimit.NewRedisSlidingWindowLimiter(redisClient, 1*time.Second, 100)
	server.Use(ratelimit.NewBuilder(limiter).Build())

	//配置session中间件
	//用于解决HTTP协议的无状态性问题
	//
	//初始化一个cookie验证的会话存储器
	//由于会话数据的保密性需求，使用cookie.NewStore创建一个会话存储实例
	//使用"secret"作为加密密钥
	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("`4A:'n1H'U:50HQx;a1p1-er1jm3\"1)b"),
	//	[]byte("3w&0yL23d@;2R1TV+iN`Jen\\C7AA266c"))
	//store, err := redis.NewStore(16, "tcp", "124.71.99.27:6379", "",
	//	[]byte("`4A:'n1H'U:50HQx;a1p1-er1jm3\"1)b"), []byte("3w&0yL23d@;2R1TV+iN`Jen\\C7AA266c"))
	//if err != nil {
	//	panic(err)
	//}
	// 配置会话中间件
	// sessions.Sessions方法创建一个名为"mysession"的会话
	// 使用上一行创建的store来存储会话数据
	// 通过server.Use将该会话中间件挂载到http服务器中
	// 以便每个请求都能进行会话管理
	//server.Use(sessions.Sessions("webook", store))

	// 此部分用于验证用户是否已登录
	//	server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup", "/users/login").Build())

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup", "/users/login").Build())
	return
}

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
