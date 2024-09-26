package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// User 直接对应数据库表结构
// 有些人叫做 entity,有些人叫做 model,有些人叫做 PO(persistent object)
type User struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`
	//全部用户唯一
	Email    string `gorm:"unique"`
	Password string

	//创建时间，毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}
