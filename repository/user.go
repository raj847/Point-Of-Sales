package repository

import (
	"context"
	"vandesar/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAdminByID(ctx context.Context, id uint) (entity.Admin, error) {
	res := entity.Admin{}
	err := r.db.WithContext(ctx).Table("admins").Where("id = ?", id).Find(&res).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return res, nil
}

func (r *userRepository) GetAdminByEmail(ctx context.Context, email string) (entity.Admin, error) {
	res := entity.Admin{}
	err := r.db.WithContext(ctx).Table("admins").Where("email = ?", email).Find(&res).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return res, nil
}

func (r *userRepository) CreateAdmin(ctx context.Context, user entity.Admin) (entity.Admin, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.Admin{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateAdmin(ctx context.Context, user entity.Admin) (entity.Admin, error) {
	err := r.db.WithContext(ctx).Table("admins").Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteAdmin(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Admin{}, id).Error
	return err
}

func (r *userRepository) GetCashierByID(ctx context.Context, id uint) (entity.Cashier, error) {
	res := entity.Cashier{}
	err := r.db.WithContext(ctx).Table("cashiers").Where("id = ?", id).Find(&res).Error
	if err != nil {
		return entity.Cashier{}, err
	}

	return res, nil
}

func (r *userRepository) GetCashierByUsername(ctx context.Context, username string) (entity.Cashier, error) {
	res := entity.Cashier{}
	err := r.db.WithContext(ctx).Table("cashiers").Where("username = ?", username).Find(&res).Error
	if err != nil {
		return entity.Cashier{}, err
	}

	return res, nil
}

func (r *userRepository) CreateCashier(ctx context.Context, user entity.Cashier) (entity.Cashier, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.Cashier{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateCashier(ctx context.Context, user entity.Cashier) (entity.Cashier, error) {
	err := r.db.WithContext(ctx).Table("cashiers").Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		return entity.Cashier{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteCashier(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Cashier{}, id).Error
	return err
}
