package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vandesar/entity"
	"vandesar/service"

	"gorm.io/gorm"
)

type ExpensesAPI struct {
	expensesService *service.ExpensesService
	userService     *service.UserService
}

func NewExpensesAPI(
	expensesService *service.ExpensesService,
	userService *service.UserService,
) *ExpensesAPI {
	return &ExpensesAPI{
		expensesService: expensesService,
		userService:     userService,
	}
}

func (p *ExpensesAPI) GetAllBeban(w http.ResponseWriter, r *http.Request) {
	adminId := r.Context().Value("id").(uint)

	beban := r.URL.Query()
	bebanID, foundBebanId := beban["beban_id"]

	if foundBebanId {
		bID, _ := strconv.Atoi(bebanID[0])
		bebanByID, err := p.expensesService.GetBebanByID(r.Context(), bID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if bebanByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error beban not found"))
			return
		}

		if bebanByID.UserID != adminId {
			WriteJSON(w, http.StatusUnauthorized, entity.NewErrorResponse("error unauthorized user id"))
		}

		WriteJSON(w, http.StatusOK, bebanByID)
		return
	}

	list, err := p.expensesService.GetBebans(r.Context(), adminId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (p *ExpensesAPI) GetAllPrive(w http.ResponseWriter, r *http.Request) {
	adminId := r.Context().Value("id").(uint)

	prive := r.URL.Query()
	priveID, foundPriveId := prive["prive_id"]

	if foundPriveId {
		prID, _ := strconv.Atoi(priveID[0])
		priveByID, err := p.expensesService.GetPriveByID(r.Context(), prID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if priveByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error prive not found"))
			return
		}

		if priveByID.UserID != adminId {
			WriteJSON(w, http.StatusUnauthorized, entity.NewErrorResponse("error unauthorized user id"))
		}

		WriteJSON(w, http.StatusOK, priveByID)
		return
	}

	list, err := p.expensesService.GetPrives(r.Context(), adminId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (p *ExpensesAPI) CreateNewBeban(w http.ResponseWriter, r *http.Request) {
	var beban entity.BebanRequest

	err := json.NewDecoder(r.Body).Decode(&beban)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid beban request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	prod, err := p.expensesService.AddBeban(r.Context(), entity.Beban{
		UserID: adminIdUint,
		// 	Listrik: beban.Listrik,
		// 	Sewa:    beban.Sewa,
		// 	Telepon: beban.Telepon,
		// 	Gaji:    beban.Gaji,
		// 	Lainnya: beban.Lainnya,
		Total: beban.Total,
		Jenis: beban.Jenis,
		Notes: beban.Notes,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":  adminIdUint,
		"beban_id": prod.ID,
		"message":  "success create new beban",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (p *ExpensesAPI) CreateNewPrive(w http.ResponseWriter, r *http.Request) {
	var prive entity.PriveRequest

	err := json.NewDecoder(r.Body).Decode(&prive)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid prive request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	prod, err := p.expensesService.AddPrive(r.Context(), entity.Prive{
		UserID: adminIdUint,
		Value:  prive.Value,
		Notes:  prive.Notes,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":  adminIdUint,
		"prive_id": prod.ID,
		"message":  "success create new prive",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (p *ExpensesAPI) DeleteBeban(w http.ResponseWriter, r *http.Request) {
	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	bebanID := r.URL.Query().Get("beban_id")
	prodID, _ := strconv.Atoi(bebanID)
	err := p.expensesService.DeleteBeban(r.Context(), prodID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":  adminIdUint,
		"beban_id": prodID,
		"message":  "success delete beban",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (p *ExpensesAPI) DeletePrive(w http.ResponseWriter, r *http.Request) {
	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	priveID := r.URL.Query().Get("prive_id")
	prodID, _ := strconv.Atoi(priveID)
	err := p.expensesService.DeletePrive(r.Context(), prodID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":  adminIdUint,
		"prive_id": prodID,
		"message":  "success delete prive",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (p *ExpensesAPI) UpdateBeban(w http.ResponseWriter, r *http.Request) {
	var beban entity.BebanRequest

	err := json.NewDecoder(r.Body).Decode(&beban)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid beban request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("beban_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid beban id"))
		return
	}

	bebans, err := p.expensesService.UpdateBeban(r.Context(), entity.Beban{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		// Listrik: beban.Listrik,
		// Sewa:    beban.Sewa,
		// Telepon: beban.Telepon,
		// Gaji:    beban.Gaji,
		// Lainnya: beban.Lainnya,
		Total:  beban.Total,
		Jenis:  beban.Jenis,
		UserID: adminIdUint,
		Notes:  beban.Notes,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":  adminIdUint,
		"beban_id": bebans.ID,
		"message":  "success update beban",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (p *ExpensesAPI) UpdatePrive(w http.ResponseWriter, r *http.Request) {
	var prive entity.PriveRequest

	err := json.NewDecoder(r.Body).Decode(&prive)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid prive request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("prive_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid prive id"))
		return
	}

	prives, err := p.expensesService.UpdatePrive(r.Context(), entity.Prive{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		UserID: adminIdUint,
		Value:  prive.Value,
		Notes:  prive.Notes,
	})
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":  adminIdUint,
		"beban_id": prives.ID,
		"message":  "success update prive",
	}

	WriteJSON(w, http.StatusOK, response)
}
