package api

import (
	"context"
	"encoding/json"
	"log"
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

func (p *TransactionAPI) GetAllTransactionsByAdmin(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(uint)

	transactionList, err := p.transactionService.ReadTransAdmin(userId)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(transactionList)
	return
}

func (p *TransactionAPI) GetAllTransactionsByCashier(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(uint)

	prod, err := p.transactionService.ReadTransCashier(userId)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(prod)
	return
}

func (p *TransactionAPI) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var product entity.TransactionReq

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("transaction_id")
	idInt, err := strconv.Atoi(id)

	products, err := p.transactionService.UpdateTrans(entity.TransactionReq{
		Debt:        product.Debt,
		Status:      product.Status,
		Money:       product.Money,
		CartList:    product.CartList,
		TotalPrice:  product.TotalPrice,
		TotalProfit: product.TotalProfit,
		Notes:       product.Notes,
	}, uint(idInt))
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    adminIdUint,
		"product_id": products.ID,
		"message":    "success update product",
	})
}

func (p *TransactionAPI) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	transId := r.URL.Query().Get("transaction_id")
	transIdInt, err := strconv.Atoi(transId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid transaction id"))
		return
	}

	prodID, _ := strconv.Atoi(transId)
	err = p.transactionService.DeleteTrans(uint(transIdInt))
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    adminIdUint,
		"product_id": prodID,
		"message":    "success delete product",
	})
}
