package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
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

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.GET("/profile", u.Profile)
	ug.POST("/edit", u.Edit)
	server.GET("/articles/list", u.lists)
}

func (u *UserHandler) lists(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
}

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
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, "注册成功")
	fmt.Printf("%v\n", req)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 这里登录成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 这里可以随便设置键值对,这是要放在 session 里面的值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: 60,
	})
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	profile, err := u.svc.Profile(ctx, userId.(int64))
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	err := sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "退出成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type userInfo struct {
		name     string
		birthday string
		profile  string
	}

	var infoReq userInfo
	err := ctx.Bind(&infoReq)
	if err != nil {
		return
	}

	if len(infoReq.name) > 10 {
		ctx.String(http.StatusOK, "昵称过长")
		return
	}
	if len(infoReq.birthday) != 10 ||
		(infoReq.birthday[4] != '-' && infoReq.birthday[7] != '-') {
		ctx.String(http.StatusOK, "日期格式不对")
	}
	if len(infoReq.profile) > 200 {
		ctx.String(http.StatusOK, "个人简介过长")
	}

	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	err = u.svc.Edit(ctx, domain.User{
		Id:       userId.(int64),
		Name:     infoReq.name,
		Birthday: infoReq.birthday,
		Profile:  infoReq.profile,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.String(http.StatusOK, "编辑成功")
}
