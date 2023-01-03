package repository

import (
	"context"
	"encoding/base64"
	"vandesar/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	res := entity.User{}
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).Find(&res).Error

	if err != nil {
		return entity.User{}, err
	}

	return res, nil // TODO: replace this
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	res := entity.User{}
	err := r.db.WithContext(ctx).Table("users").Where("email = ?", email).Find(&res).Error

	if err != nil {
		return entity.User{}, err
	}

	return res, nil // TODO: replace this
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}
