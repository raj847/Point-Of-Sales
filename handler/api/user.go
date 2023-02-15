package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if adminReq.Email == "" || adminReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}

	eUser, err := u.userService.LoginAdmin(r.Context(), adminReq)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			w.WriteHeader(404)
			w.Write([]byte(err.Error()))
			return
		} else if strings.Contains(err.Error(), "email not found") {
			w.WriteHeader(404)
			w.Write([]byte(err.Error()))
			return
		} else if strings.Contains(err.Error(), "password not match") {
			w.WriteHeader(404)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &entity.Claims{
		UserID:  eUser.ID,
		AdminID: eUser.ID,
		Role:    "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("rahasia-perusahaan"))

	//set cookies
	expiresAt := time.Now().Add(5 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Path:    "/",
		Value:   tokenString,
		Expires: expiresAt,
	})

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": int(eUser.ID),
		"role":    "admin",
		"message": "login success",
	})
}

func (u *UserAPI) AdminRegister(w http.ResponseWriter, r *http.Request) {
	var user entity.AdminRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.ShopName == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	eUser, err := u.userService.RegisterAdmin(r.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("email already exists"))
			return
		} else if strings.Contains(err.Error(), "format email invalid") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("format email invalid"))
			return
		} else if strings.Contains(err.Error(), "your domain not found") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("your domain not found"))
			return
		} else if strings.Contains(err.Error(), "password is not valid") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("password is not valid"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": eUser.ID,
		"message": "register success",
	})
}

func (u *UserAPI) CashierLogin(w http.ResponseWriter, r *http.Request) {
	var cashierReq entity.CashierLogin

	err := json.NewDecoder(r.Body).Decode(&cashierReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if cashierReq.Username == "" || cashierReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}

	eUser, err := u.userService.LoginCashier(r.Context(), cashierReq)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &entity.Claims{
		UserID:  eUser.ID,
		AdminID: eUser.AdminID,
		Role:    "cashier",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("rahasia-perusahaan"))

	//set cookies
	expiresAt := time.Now().Add(5 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Path:    "/",
		Value:   tokenString,
		Expires: expiresAt,
	})

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": int(eUser.ID),
		"message": "login success",
	})
}

func (u *UserAPI) CashierRegister(w http.ResponseWriter, r *http.Request) {
	var user entity.CashierRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	id := r.Context().Value("id").(string)

	adminId := strings.Split(id, "|")[1] // admin id
	adminIdUint, err := strconv.Atoi(adminId)
	user.AdminID = uint(adminIdUint)

	eUser, err := u.userService.RegisterCashier(r.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "password is not valid") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("password is not valid"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": eUser.ID,
		"role":    "cashier",
		"message": "register success",
	})
}

func (u *UserAPI) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
}
