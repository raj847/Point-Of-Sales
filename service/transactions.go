package service

import (
	"vandesar/entity"
	"vandesar/repository"
)

type TransactionService interface {
	AddTrans(trans entity.TransactionReq) []error
	UpdateTrans(trans entity.Transaction) (entity.Transaction, error)
	DeleteTrans(id int) error
	ReadTrans() ([]entity.TransactionReq, error)
}

type transactionService struct {
	transRepo repository.TransactionRepository
}

func NewTransactionService(transRepo repository.TransactionRepository) TransactionService {
	return &transactionService{transRepo}
}

func (t *transactionService) AddTrans(trans entity.TransactionReq) []error {
	err := t.transRepo.AddTrans(trans)
	if err != nil {
		return err
	}
	return nil
}
func (t *transactionService) UpdateTrans(trans entity.Transaction) (entity.Transaction, error) {
	err := t.transRepo.UpdateTrans(trans)
	if err != nil {
		return entity.Transaction{}, err
	}
	return trans, nil
}
func (t *transactionService) DeleteTrans(id int) error {
	return t.transRepo.DeleteTrans(id)
}
func (t *transactionService) ReadTrans() ([]entity.TransactionReq, error) {
	return t.transRepo.ReadTrans()
}
