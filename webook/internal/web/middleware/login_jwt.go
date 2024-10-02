package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/webook/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(paths ...string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, paths...)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要登录校验
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}

		// 现在用 JWT 来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//segs := strings.SplitN(tokenHeader, " ", 2)
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		// ParseWithClaims 里面一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("\"`4A:'n1H'U:50HQx;a1p1-er1jm3\\\"1)b\""), nil
		})
		if err != nil || (token == nil || !token.Valid) || claims.Uid == 0 {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//每10秒刷新一次
		now := time.Now()
		//这里检查 JWT 令牌的过期时间（claims.ExpiresAt）和当前时间的差值，判断是否少于 50 秒。
		//如果剩余时间少于 50 秒，则意味着令牌快要过期，代码会进入 if 块，执行续签操作。
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("\"`4A:'n1H'U:50HQx;a1p1-er1jm3\\\"1"))
			if err != nil {
				//记录日志
				log.Println("jwt 缓存失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
		//ctx.Set("UserId", claims.Uid)
	}
}
