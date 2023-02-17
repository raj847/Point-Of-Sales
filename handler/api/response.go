package api

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	response := SuccessResponse{
		Code:   code,
		Status: http.StatusText(code),
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
