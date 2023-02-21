package entity

import (
	"time"

	"gorm.io/gorm"
)

type Rekap struct {
	gorm.Model
	AdminID uint `json:"admin_id"`

	TotalPrice float64 `json:"total_price"`
	TotalProfit float64 `json:"total_profit"`
	TotalDebt float64 `json:"total_debt"`
	TotalPeopleDebt int `json:"total_people_debt"`

	LinkPdf string `json:"link_pdf"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
}

type RekapPerDay struct {
	gorm.Model
	AdminID uint `json:"admin_id"`

	TotalPrice float64 `json:"total_price"`
	TotalProfit float64 `json:"total_profit"`
	TotalDebt float64 `json:"total_debt"`
	TotalPeopleDebt int `json:"total_people_debt"`

	LinkPdf string `json:"link_pdf"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
}
