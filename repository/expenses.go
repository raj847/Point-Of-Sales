package repository

import (
	"context"
	"vandesar/entity"

	"gorm.io/gorm"
)

type ExpensesRepository struct {
	db *gorm.DB
}

func NewExpensesRepository(db *gorm.DB) *ExpensesRepository {
	return &ExpensesRepository{db}
}

func (p *ExpensesRepository) GetBebanByUserId(ctx context.Context, userId uint) ([]entity.Beban, error) {
	var bebansResult []entity.Beban

	bebans, err := p.db.
		WithContext(ctx).
		Table("bebans").
		Select("*").
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Rows()
	if err != nil {
		return []entity.Beban{}, err
	}
	defer bebans.Close()

	for bebans.Next() {
		p.db.ScanRows(bebans, &bebansResult)
	}

	return bebansResult, nil
}

func (p *ExpensesRepository) GetPriveByUserId(ctx context.Context, userId uint) ([]entity.Prive, error) {
	var privesResult []entity.Prive

	prives, err := p.db.
		WithContext(ctx).
		Table("prives").
		Select("*").
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Rows()
	if err != nil {
		return []entity.Prive{}, err
	}
	defer prives.Close()

	for prives.Next() {
		p.db.ScanRows(prives, &privesResult)
	}

	return privesResult, nil
}

func (p *ExpensesRepository) AddBeban(ctx context.Context, bebans entity.Beban) error {
	err := p.db.
		WithContext(ctx).
		Create(&bebans).Error
	return err
}

func (p *ExpensesRepository) AddPrive(ctx context.Context, prives entity.Prive) error {
	err := p.db.
		WithContext(ctx).
		Create(&prives).Error
	return err
}

func (p *ExpensesRepository) GetBebanByID(ctx context.Context, id int) (entity.Beban, error) {
	var bebansResult entity.Beban

	err := p.db.
		WithContext(ctx).
		Table("bebans").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&bebansResult).Error
	if err != nil {
		return entity.Beban{}, err
	}

	return bebansResult, nil
}

func (p *ExpensesRepository) GetPriveByID(ctx context.Context, id int) (entity.Prive, error) {
	var privesResult entity.Prive

	err := p.db.
		WithContext(ctx).
		Table("prives").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&privesResult).Error
	if err != nil {
		return entity.Prive{}, err
	}

	return privesResult, nil
}

func (p *ExpensesRepository) DeleteBeban(ctx context.Context, id int) error {
	err := p.db.
		WithContext(ctx).
		Delete(&entity.Beban{}, id).Error
	return err
}

func (p *ExpensesRepository) DeletePrive(ctx context.Context, id int) error {
	err := p.db.
		WithContext(ctx).
		Delete(&entity.Prive{}, id).Error
	return err
}

func (p *ExpensesRepository) UpdateBeban(ctx context.Context, bebans entity.Beban) error {
	err := p.db.
		WithContext(ctx).
		Table("bebans").
		Where("id = ?", bebans.ID).
		Updates(&bebans).Error
	return err
}

func (p *ExpensesRepository) UpdatePrive(ctx context.Context, prives entity.Prive) error {
	err := p.db.
		WithContext(ctx).
		Table("prives").
		Where("id = ?", prives.ID).
		Updates(&prives).Error
	return err
}
