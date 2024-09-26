package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(124.71.99.27:13306)/webook"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		// 只在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化出错，应用就不要启动了
		panic(err)
	}

	if err := dao.InitTable(db); err != nil {
		panic(err)
	}

	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)

	server := web.RegisterRoutes(u)
	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
