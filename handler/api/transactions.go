package api

import (
	"net/http"
	"vandesar/service"
)

type TransactionAPI interface {
	GetTransaction(w http.ResponseWriter, r *http.Request)
	CreateNewTransaction(w http.ResponseWriter, r *http.Request)
	DeleteTransaction(w http.ResponseWriter, r *http.Request)
	UpdateTransaction(w http.ResponseWriter, r *http.Request)
}

type transactionAPI struct {
	transactionService service.TransactionService
}

func NewTransactionAPI(transactionService service.TransactionService) *transactionAPI {
	return &transactionAPI{transactionService}
}

func (p *transactionAPI) GetTransaction(w http.ResponseWriter, r *http.Request) {

}
