package entity

import (
	"gorm.io/gorm"
)

type Prods struct {
	Id         int
	ProductID  int
	Quantity   int
	TotalPrice float64
}

type Transaction struct {
	gorm.Model
	UserID     int     `json:"user_id"`
	Debt       float64 `json:"debt"`
	Status     string  `json:"status"`
	Money      float64 `json:"money"`
	Products   []Prods `json:listproduct` //ini isinya produk yang dibeli + quantitynya + harga total per product dari frontend
	TotalHarga float64 `json:"totalduit`
	Notes      string  `json:"notes"`
}
