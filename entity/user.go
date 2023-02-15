package entity

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ShopName string `json:"shop_name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	Role     string `json:"role" gorm:"type:varchar(50);not null" validate:"required"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
}

type Cashier struct {
	gorm.Model
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username" gorm:"type:varchar(255);not null" validate:"required"`
	Role     string `json:"role" gorm:"type:varchar(50);not null" validate:"required"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
}

type AdminLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminRegister struct {
	ShopName string `json:"shop_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"-"`
	Password string `json:"password" binding:"required"`
}

type CashierLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CashierRegister struct {
	AdminID  uint   `json:"admin_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"-"`
	Password string `json:"password" binding:"required"`
}
