package entity

import (
	"gorm.io/gorm"
)

type Prods struct {
	ProductID  uint    `json:"product_id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Price      float64 `json:"price"`
	Quantity   uint    `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type TransactionReq struct {
	UserID      uint    `json:"user_id"`
	Debt        float64 `json:"debt"`
	Status      string  `json:"status"`
	Money       float64 `json:"money"`
	CartList    []Prods `json:"cart_list"`
	TotalPrice  float64 `json:"total_price"`
	TotalProfit float64 `json:"total_profit"`
	Notes       string  `json:"notes"`
}

type Transaction struct {
	gorm.Model
	UserID      uint    `json:"user_id"` // cashier_id
	Debt        float64 `json:"debt"`
	Status      string  `json:"status"`
	Money       float64 `json:"money"`
	CartList    []byte  `json:"cart_list"`
	TotalPrice  float64 `json:"total_price"`
	Notes       string  `json:"notes"`
	TotalProfit float64 `json:"total_profit"`
}
