package entity

import (
	"gorm.io/gorm"
)

type Beban struct {
	gorm.Model

	// Listrik float64 `json:"listrik"`
	// Sewa    float64 `json:"sewa"`
	// Telepon float64 `json:"telepon"`
	// Gaji    float64 `json:"gaji"`
	Total  float64 `json:"total"`
	UserID uint    `json:"user_id"`
	Jenis  string  `json:"jenis"`
}

type BebanRequest struct {
	// Listrik float64 `json:"listrik"`
	// Sewa    float64 `json:"sewa"`
	// Telepon float64 `json:"telepon"`
	// Gaji    float64 `json:"gaji"`
	// Lainnya float64 `json:"lainnya"`
	Total float64 `json:"total"`
	Jenis string  `json:"jenis"`
}

type Prive struct {
	gorm.Model

	Value  float64 `json:"value"`
	Notes  string  `json:"notes"`
	UserID uint    `json:"user_id"`
}

type PriveRequest struct {
	Value float64 `json:"value"`
	Notes string  `json:"notes"`
}
