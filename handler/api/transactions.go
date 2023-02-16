package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vandesar/entity"
	"vandesar/service"
)

type TransactionAPI struct {
	transactionService *service.TransactionService
}

func NewTransactionAPI(transactionService *service.TransactionService) *TransactionAPI {
	return &TransactionAPI{transactionService}
}

func (p *TransactionAPI) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var product entity.TransactionReq

	ctx := r.Context()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid transaction request"))
		return
	}

	cashierIdUint := r.Context().Value("id").(uint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	errs := p.transactionService.AddTrans(ctxWithTimeout, entity.TransactionReq{
		UserID:      cashierIdUint,
		Debt:        product.Debt,
		Status:      product.Status,
		Money:       product.Money,
		CartList:    product.CartList,
		TotalPrice:  product.TotalPrice,
		TotalProfit: product.TotalProfit,
		Notes:       product.Notes,
	})
	if len(errs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": cashierIdUint,
		"message": "success create new transaction",
	})
}

func (p *TransactionAPI) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	cashierIdUint := r.Context().Value("id").(uint)

	userType := r.URL.Query()

	adminId, foundAdminId := userType["admin"]
	adminIdUint, _ := strconv.Atoi(adminId[0])

	if foundAdminId {
		transactionList, err := p.transactionService.ReadTransAdmin(adminIdUint)
		if err != nil {
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(transactionList)
		return
	}

	prod, err := p.transactionService.ReadTransCashier(cashierIdUint)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(prod)
	return
}
