package repository

import (
	"context"
	"vandesar/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetAllAdmins() ([]entity.Admin, error) {
	var admins []entity.Admin
	err := r.db.Table("admins").Find(&admins).Error
	if err != nil {
		return []entity.Admin{}, err
	}

	return admins, nil
}

func (r *UserRepository) GetAdminByID(ctx context.Context, id uint) (entity.Admin, error) {
	var adminResult entity.Admin
	err := r.db.WithContext(ctx).Table("admins").Where("id = ?", id).Find(&adminResult).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return adminResult, nil
}

func (r *UserRepository) ChangeAdminPassword(ctx context.Context, id uint, password string) error {
	err := r.db.WithContext(ctx).Table("admins").Where("id = ?", id).Update("password", password).Error
	return err
}

func (r *UserRepository) GetAdminByEmail(ctx context.Context, email string) (entity.Admin, error) {
	var adminResult entity.Admin
	err := r.db.WithContext(ctx).Table("admins").Where("email = ?", email).Find(&adminResult).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return adminResult, nil
}

func (r *UserRepository) CreateAdmin(ctx context.Context, user entity.Admin) (entity.Admin, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdateAdmin(ctx context.Context, user entity.Admin) (entity.Admin, error) {
	err := r.db.WithContext(ctx).Table("admins").Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return user, nil
}

func (r *UserRepository) DeleteAdmin(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Admin{}, id).Error
	return err
}

func (r *UserRepository) GetCashierByID(ctx context.Context, id uint) (entity.Cashier, error) {
	var res entity.Cashier
	err := r.db.WithContext(ctx).Table("cashiers").Where("id = ?", id).Find(&res).Error
	if err != nil {
		return entity.Cashier{}, err
	}

	return res, nil
}

func (r *UserRepository) GetCashierByUsername(ctx context.Context, username string) (entity.Cashier, error) {
	var res entity.Cashier
	err := r.db.WithContext(ctx).Table("cashiers").Where("username = ?", username).Find(&res).Error
	if err != nil {
		return entity.Cashier{}, err
	}

	return res, nil
}

func (r *UserRepository) CreateCashier(ctx context.Context, user entity.Cashier) (entity.Cashier, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.Cashier{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateCashier(ctx context.Context, user entity.Cashier) (entity.Cashier, error) {
	err := r.db.WithContext(ctx).Table("cashiers").Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		return entity.Cashier{}, err
	}

	return user, nil
}

func (r *UserRepository) DeleteCashier(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Cashier{}, id).Error
	return err
}

func (r *UserRepository) CheckTokenAdmin(token entity.CheckTokenAdmin) (error) {
	err := r.db.Create(&token).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CheckTokenCashier(token entity.CheckTokenCashier) (error) {
	err := r.db.Create(&token).Error
	if err != nil {
		return err
	}

	return nil
}