package entity

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ShopName string `json:"shop_name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
	Role     string `json:"role" gorm:"type:varchar(50);not null" validate:"required"`
	// PhotoURL string `json:"photo_url" gorm:"type:varchar(255)"`
}

type Cashier struct {
	gorm.Model
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username" gorm:"type:varchar(255);not null" validate:"required"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
	Role     string `json:"role" gorm:"type:varchar(50);not null" validate:"required"`
	// Online   bool   `json:"online" gorm:"type:boolean`
}

type AdminLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type AdminChangePassword struct {
	AdminID     uint   `json:"-"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type AdminRegister struct {
	ShopName string `json:"shop_name" binding:"required"`

	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"-"`
	// PhotoURL string `json:"photo_url"`
}

type CashierLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	// Online   bool   `json:"-"`
}

type CashierRegister struct {
	AdminID uint `json:"admin_id" binding:"required"`

	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"-"`
}

// type CheckTokenAdmin struct {
// 	TokenInput string `json:"token_input"`
// 	AdminID    uint   `json:"admin_id"`
// }

// type CheckTokenCashier struct {
// 	TokenInput string `json:"token_input"`
// 	CashierId    uint   `json:"cashier_id"`
// }

type CheckToken struct {
	TokenInput string `json:"token_input"`
	UserId     uint   `json:"user_id"`
}
