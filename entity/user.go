package entity

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	ShopName  string    `json:"shopname" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	Role      string    `json:"role" gorm:"type:varchar(50);not null" validate:"required"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cashier struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	AdminID   int       `json:"admin_id"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	Role      string    `json:"role" gorm:"type:varchar(50);not null" validate:"required"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	ShopName string `json:"shopname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}
