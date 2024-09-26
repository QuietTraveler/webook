package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/webook/internal/domain"
	"webook/webook/internal/service"
)

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[\\w](?:[\\w-]*[\\w])?"
		passwordRegexPattern = "^(?=.*\\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{8,10}$"
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

//func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
//	userGroup := server.Group("/users")
//	{
//		userGroup.POST("/signup", u.SignUp)
//		userGroup.POST("/login", u.Login)
//		userGroup.POST("/edit", u.Edit)
//		userGroup.GET("/profile", u.Profile)
//	}
//
//}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// 解析请求
	//Bind 方法会根据 Context-Type 来解析你的数据到 req 里面
	//解析错了，就会放回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	//ok, err := regexp.Match(emailRegexPattern, []byte(req.Email))
	if err != nil {
		// 记录日志
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusBadRequest, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusBadRequest, "两次输入的密码不一致")
		return
	}

	//	ok, err = regexp.Match(passwordRegexPattern, []byte(req.Password))
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		//记录日志
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		// 记录日志
		ctx.String(http.StatusBadRequest, "密码必须包含大小写字母和数字的组合，不能使用特殊字符，长度在8-10之间")
		return
	}

	// 调用 service 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, "注册成功")
	fmt.Printf("%v\n", req)
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}
