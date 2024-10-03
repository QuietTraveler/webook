package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 在这里操作缓存
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//先从 cache 里面找
	// 再从 dao 里面找
	//找到了回写 cache
	user, err := r.dao.FindById(ctx, id)
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Profile:  user.Profile,
		Birthday: user.Birthday,
	}, err
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, dao.User{
		Id:       u.Id,
		Name:     u.Name,
		Profile:  u.Profile,
		Birthday: u.Birthday,
	})
}
