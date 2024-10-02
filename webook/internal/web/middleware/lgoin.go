package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) IgnorePaths(paths ...string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, paths...)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用Go的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要登录校验
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}

		////不需要登录校验
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}

		// 尝试从上下文中获取默认的会话
		sess := sessions.Default(ctx)
		// 尝试从会话中获取用户ID
		id := sess.Get("userId")
		// 如果用户ID为nil，表示用户未登录
		if id == nil {
			// 终止请求，返回未授权状态码
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		//刚登录，还没有刷新
		if updateTime == nil {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}

		//updateTime 是有的
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}
	}
}
