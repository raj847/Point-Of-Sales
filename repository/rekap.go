package repository

import (
	"vandesar/entity"

	"gorm.io/gorm"
)

type RekapRepository struct {
	db *gorm.DB
}

func NewRekapRepository(db *gorm.DB) *RekapRepository {
	return &RekapRepository{db}
}

func (p *RekapRepository) ListRekap(adminId uint) ([]entity.Rekap, error) {
	var rekap []entity.Rekap

	err := p.db.
	Debug().
		Table("rekaps").
		Select("*").
		Where("admin_id = ?", adminId).
		Where("deleted_at IS NULL").
		Find(&rekap).Error

	return rekap, err
}

func (p *RekapRepository) ListRekapPerDays(adminId uint) ([]entity.Rekap, error) {
	var rekap []entity.Rekap

	err := p.db.
	Debug().
		Table("rekap_per_days").
		Select("*").
		Where("admin_id = ?", adminId).
		Where("deleted_at IS NULL").
		Find(&rekap).Error

	return rekap, err
}

func (p *RekapRepository) AddRekap(rekap entity.Rekap) error {
	err := p.db.
		Create(&rekap).Error
	return err
}

func (p *RekapRepository) AddRekapPerDay(rekap entity.RekapPerDay) error {
	err := p.db.
		Create(&rekap).Error
	return err
}
