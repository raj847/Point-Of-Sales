package service

import (
	"context"
	"errors"
	"fmt"
	"vandesar/entity"
	"vandesar/repository"
)

type ExpensesService struct {
	exRepo *repository.ExpensesRepository
}

func NewExpensesService(exRepo *repository.ExpensesRepository) *ExpensesService {
	return &ExpensesService{
		exRepo: exRepo,
	}
}

func (s *ExpensesService) GetBebans(ctx context.Context, adminId uint) ([]entity.Beban, error) {
	return s.exRepo.GetBebanByUserId(ctx, adminId)
}

func (s *ExpensesService) GetPrives(ctx context.Context, adminId uint) ([]entity.Prive, error) {
	return s.exRepo.GetPriveByUserId(ctx, adminId)
}

func (s *ExpensesService) AddBeban(ctx context.Context, beban entity.Beban) (entity.Beban, error) {
	err := s.exRepo.AddBeban(ctx, beban)
	if err != nil {
		return entity.Beban{}, err
	}
	return beban, nil
}

func (s *ExpensesService) AddPrive(ctx context.Context, prive entity.Prive) (entity.Prive, error) {
	err := s.exRepo.AddPrive(ctx, prive)
	if err != nil {
		return entity.Prive{}, err
	}
	return prive, nil
}

func (s *ExpensesService) GetBebanByID(ctx context.Context, id int) (entity.Beban, error) {
	return s.exRepo.GetBebanByID(ctx, id)
}

func (s *ExpensesService) GetPriveByID(ctx context.Context, id int) (entity.Prive, error) {
	return s.exRepo.GetPriveByID(ctx, id)
}

func (s *ExpensesService) UpdateBeban(ctx context.Context, beban entity.Beban) (entity.Beban, error) {
	existingBeban, err := s.GetBebanByID(ctx, int(beban.ID))
	if err != nil {
		return entity.Beban{}, err
	}

	if existingBeban.UserID != beban.UserID {
		return entity.Beban{}, errors.New("you are not allowed to update this beban")
	}

	err = s.exRepo.UpdateBeban(ctx, beban)
	if err != nil {
		return entity.Beban{}, err
	}

	return beban, nil
}

func (s *ExpensesService) UpdatePrive(ctx context.Context, prive entity.Prive) (entity.Prive, error) {
	existingPrive, err := s.GetBebanByID(ctx, int(prive.ID))
	if err != nil {
		return entity.Prive{}, err
	}

	if existingPrive.UserID != prive.UserID {
		fmt.Println(existingPrive.UserID, prive.UserID)
		return entity.Prive{}, errors.New("you are not allowed to update this prive")
	}

	err = s.exRepo.UpdatePrive(ctx, prive)
	if err != nil {
		return entity.Prive{}, err
	}

	return prive, nil
}

func (s *ExpensesService) DeleteBeban(ctx context.Context, id int) error {
	return s.exRepo.DeleteBeban(ctx, id)
}

func (s *ExpensesService) DeletePrive(ctx context.Context, id int) error {
	return s.exRepo.DeletePrive(ctx, id)
}
