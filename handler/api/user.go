package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vandesar/entity"
	"vandesar/service"

	"github.com/dgrijalva/jwt-go"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	if user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}
	eUser, err := u.userService.Login(r.Context(), &entity.User{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	//set jwt
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &entity.Claims{
		UserID: eUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("rahasia-perusahaan"))

	//set cookies
	expiresAt := time.Now().Add(5 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Path:    "/",
		Value:   tokenString,
		Expires: expiresAt,
	})

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": strconv.Itoa(eUser), "message": "login success"})
	// TODO: answer here
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

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
	fmt.Println(user)
	eUser, err := u.userService.Register(r.Context(), &entity.User{
		ShopName: user.ShopName,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("email already exists"))
			return
		} else if strings.Contains(err.Error(), "format email invalid") {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("format email invalid"))
			return
		} else if strings.Contains(err.Error(), "your domain not found") {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("your domain not found"))
			return
		} else if strings.Contains(err.Error(), "password is not valid") {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("password is not valid"))
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": eUser.ID,
		"message": "register success"})
	// TODO: answer here
}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
