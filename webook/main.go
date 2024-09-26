package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/webook/internal/repository"
	"webook/webook/internal/repository/dao"
	"webook/webook/internal/service"
	"webook/webook/internal/web"
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

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}

	u := initUser(db)
	server := web.RegisterRoutes(u)
	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
