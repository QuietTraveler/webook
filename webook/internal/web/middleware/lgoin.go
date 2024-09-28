package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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
	return func(ctx *gin.Context) {
		//不需要登录校验
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}
		//不需要登录校验
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}
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
	}
}
