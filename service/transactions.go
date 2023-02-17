package service

import (
	"context"
	"vandesar/entity"
	"vandesar/repository"
)

type TransactionService struct {
	transRepo *repository.TransactionRepository
}

func NewTransactionService(transRepo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{transRepo}
}

func (t *TransactionService) AddTrans(ctx context.Context, trans entity.TransactionReq) []error {
	err := t.transRepo.AddTrans(ctx, trans)
	if err != nil {
		return err
	}
	return nil
}
func (t *TransactionService) UpdateTrans(trans entity.TransactionReq, tranId uint) (entity.Transaction, error) {
	result, err := t.transRepo.UpdateTrans(trans, tranId)
	if err != nil {
		return entity.Transaction{}, err
	}

	return result, nil
}
func (t *TransactionService) DeleteTrans(id uint) error {
	return t.transRepo.DeleteTrans(id)
}

func (t *TransactionService) ReadTransCashier(id uint) ([]entity.TransactionReq, error) {
	return t.transRepo.ReadTransByCashier(id)
}

func (t *TransactionService) ReadTransAdmin(adminId uint) ([]entity.TransactionReq, error) {
	return t.transRepo.ReadTransByAdmin(adminId)
}
