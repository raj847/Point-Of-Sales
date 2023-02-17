package entity

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Code   string  `json:"code"`
	Name   string  `gorm:"type:varchar(255);unique_index" json:"name"`
	Price  float64 `json:"price"`
	Stock  uint    `json:"stock"`

	// UserID reference to adminId in Admins table
	// because just admin can create product
	UserID uint    `json:"user_id"`

	Modal  float64 `json:"modal"`
}

type ProductRequest struct {
	Code  string  `json:"code"`
	Name  string  `gorm:"type:varchar(255);unique_index" json:"name"`
	Price float64 `json:"price"`
	Stock uint    `json:"stock"`
	Modal float64 `json:"modal"`
}
