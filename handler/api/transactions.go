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

	parentCtx := r.Context()
	ctxWithTimeout, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	cashierIdUint := r.Context().Value("id").(uint)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
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
		WriteJSON(w, http.StatusBadRequest, errs)
		return
	}

	response := map[string]any{
		"user_id": cashierIdUint,
		"message": "success create new transaction",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (p *TransactionAPI) GetAllTransactionsByAdmin(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(uint)

	transactionList, err := p.transactionService.ReadTransAdmin(userId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, transactionList)
}

func (p *TransactionAPI) GetAllTransactionsByAdminDebt(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(uint)

	transactionList, err := p.transactionService.ReadTransAdminDebt(userId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, transactionList)
}

func (p *TransactionAPI) GetAllTransactionsByCashier(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(uint)

	prod, err := p.transactionService.ReadTransCashier(userId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, prod)
}

func (p *TransactionAPI) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var product entity.TransactionReq

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("transaction_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid transaction id"))
		return
	}

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
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":        adminIdUint,
		"transaction_id": products.ID,
		"message":        "success update transaction",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (p *TransactionAPI) UpdateTransactionDebt(w http.ResponseWriter, r *http.Request) {
	var product entity.TransactionReq

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("transaction_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid transaction id"))
		return
	}

	products, err := p.transactionService.UpdateTransDebt(entity.UpdateTrans{
		Debt:   product.Debt,
		Status: product.Status,
		Money:  product.Money,
	}, uint(idInt))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":        adminIdUint,
		"transaction_id": products.ID,
		"message":        "success update debt transaction",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (p *TransactionAPI) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	transId := r.URL.Query().Get("transaction_id")
	transIdInt, err := strconv.Atoi(transId)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid transaction id"))
		return
	}

	prodID, _ := strconv.Atoi(transId)
	err = p.transactionService.DeleteTrans(uint(transIdInt))
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":    adminIdUint,
		"product_id": prodID,
		"message":    "success delete product",
	}

	WriteJSON(w, http.StatusOK, response)
}
