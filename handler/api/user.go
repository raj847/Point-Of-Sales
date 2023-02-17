package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"vandesar/entity"
	"vandesar/service"

	"github.com/dgrijalva/jwt-go"
)

type UserAPI struct {
	userService *service.UserService
}

func NewUserAPI(userService *service.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

func (u *UserAPI) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var adminReq entity.AdminLogin

	err := json.NewDecoder(r.Body).Decode(&adminReq)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if adminReq.Email == "" || adminReq.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("email or password is empty"))
		return
	}

	eUser, err := u.userService.LoginAdmin(r.Context(), adminReq)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserNotFound) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserPasswordDontMatch) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)

	claims := entity.Claims{
		UserID:  eUser.ID,
		AdminID: eUser.ID,
		Role:    "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString([]byte("rahasia-perusahaan"))

	expiresAt := time.Now().Add(5 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Path:    "/",
		Value:   tokenString,
		Expires: expiresAt,
	})

	response := map[string]any{
		"user_id": int(eUser.ID),
		"role":    "admin",
		"message": "login success",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (u *UserAPI) AdminRegister(w http.ResponseWriter, r *http.Request) {
	var user entity.AdminRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.ShopName == "" || user.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("register data is empty"))
		return
	}

	eUser, err := u.userService.RegisterAdmin(r.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrEmailInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrPasswordInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id": eUser.ID,
		"message": "register success",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (u *UserAPI) CashierLogin(w http.ResponseWriter, r *http.Request) {
	var cashierReq entity.CashierLogin

	err := json.NewDecoder(r.Body).Decode(&cashierReq)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if cashierReq.Username == "" || cashierReq.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("username or password is empty"))
		return
	}

	eUser, err := u.userService.LoginCashier(r.Context(), cashierReq)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserPasswordDontMatch) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := entity.Claims{
		UserID:  eUser.ID,
		AdminID: eUser.AdminID,
		Role:    "cashier",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString([]byte("rahasia-perusahaan"))

	expiresAt := time.Now().Add(5 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Path:    "/",
		Value:   tokenString,
		Expires: expiresAt,
	})

	response := map[string]any{
		"user_id": int(eUser.ID),
		"role":    "cashier",
		"message": "login success",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (u *UserAPI) CashierRegister(w http.ResponseWriter, r *http.Request) {
	var user entity.CashierRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Username == "" || user.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("register data is empty"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	user.AdminID = uint(adminIdUint)

	eUser, err := u.userService.RegisterCashier(r.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrPasswordInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserPasswordDontMatch) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id": eUser.ID,
		"role":   "cashier",
		"message": "register success",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (u *UserAPI) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
}
