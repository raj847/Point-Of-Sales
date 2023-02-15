package entity

import "time"

type Product struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Code      string    `json:"code"`
	Name      string    `gorm:"type:varchar(255);unique_index" json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	UserID    int       `json:"user_id"`
	Modal     float64   `json:"modal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductRequest struct {
	ID    int     `json:"id"`
	Code  string  `json:"code"`
	Name  string  `gorm:"type:varchar(255);unique_index" json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	Modal float64 `json:"modal"`
}
