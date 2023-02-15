package entity

import (
	"gorm.io/gorm"
)

type Prods struct {
	ProductID  int
	Quantity   int
	TotalPrice float64
}

type TransactionReq struct {
	UserID     int     `json:"user_id"`
	Debt       float64 `json:"debt"`
	Status     string  `json:"status"`
	Money      float64 `json:"money"`
	Products   []Prods `json:cartlist` //ini isinya produk yang dibeli + quantitynya + harga total per product dari frontend
	TotalHarga float64 `json:"totalduit`
	TotalLaba  float64 `json:"totalharga"`
	Notes      string  `json:"notes"`
}

type Transaction struct {
	gorm.Model
	UserID     int     `json:"user_id"`
	Debt       float64 `json:"debt"`
	Status     string  `json:"status"`
	Money      float64 `json:"money"`
	Products   []byte  `json:cartlist` //ini isinya produk yang dibeli + quantitynya + harga total per product dari frontend
	TotalHarga float64 `json:"totalduit`
	Notes      string  `json:"notes"`
	TotalLaba  float64 `json:"totalharga"`
}
