package entity

import (
	"gorm.io/gorm"
)

type Prods struct {
	ProductID  uint    `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Price      float64 `json:"price"`
	Quantity   uint    `json:"qty"`
	TotalPrice float64 `json:"totalPrice"`
	Modal      float32 `json:"modal"`
}

type TransactionReq struct {
	UserID      uint    `json:"cashierId"`
	Debt        float64 `json:"change"`
	Status      string  `json:"status"`
	Money       float64 `json:"pay"`
	CartList    []Prods `json:"cartList"`
	TotalPrice  float64 `json:"total"`
	TotalProfit float64 `json:"totalProfit"`
	Notes       string  `json:"notes"`
}

type Transaction struct {
	gorm.Model
	UserID      uint    `json:"cashierId"` // cashier_id
	Debt        float64 `json:"change"`
	Status      string  `json:"status"`
	Money       float64 `json:"pay"`
	CartList    []byte  `json:"cartList"`
	TotalPrice  float64 `json:"total"`
	Notes       string  `json:"notes"`
	TotalProfit float64 `json:"totalProfit"`
}

type UpdateTrans struct {
	Debt   float64 `json:"change"`
	Status string  `json:"status"`
	Money  float64 `json:"pay"`
}
