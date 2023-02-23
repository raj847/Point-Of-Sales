package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"vandesar/entity"
	"vandesar/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/minio/minio-go/v7"
)

type UserAPI struct {
	userService *service.UserService
	minioClient *minio.Client
}

func NewUserAPI(
	userService *service.UserService,
	minioClient *minio.Client,
) *UserAPI {
	return &UserAPI{
		userService: userService,
		minioClient: minioClient,
	}
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
		"user_id":     int(eUser.ID),
		"role":        "admin",
		"nama":        eUser.ShopName,
		"message":     "login success",
		"tokenCookie": tokenString,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (u *UserAPI) ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
	var changeAdminPassReq entity.AdminChangePassword

	err := json.NewDecoder(r.Body).Decode(&changeAdminPassReq)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if changeAdminPassReq.OldPassword == "" || changeAdminPassReq.NewPassword == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("email or password is empty"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	changeAdminPassReq.AdminID = adminIdUint

	eUser, err := u.userService.ChangeAdminPassword(r.Context(), changeAdminPassReq)
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
		"user_id":     int(eUser.ID),
		"role":        "admin",
		"message":     "login success",
		"tokenCookie": tokenString,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (u *UserAPI) AdminRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get form file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}

	_, err = u.minioClient.PutObject(r.Context(), "rajendra", header.Filename, file, header.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	fileName := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header.Filename)

	// get form value from request

	email := r.FormValue("email")
	password := r.FormValue("password")
	role := "admin"
	shopName := r.FormValue("shop_name")

	user := entity.AdminRegister{
		Email:    email,
		Password: password,
		Role:     role,
		ShopName: shopName,
		PhotoURL: fileName,
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
		"user_id":  int(eUser.ID),
		"role":     "cashier",
		"nama":     eUser.Username,
		"tokentod": tokenString,
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
		"role":    "cashier",
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

func (u *UserAPI) CheckToken(w http.ResponseWriter, r *http.Request) {
	var token entity.CheckToken

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if token.TokenInput == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("token nya cok"))
		return
	}

	// c, _ := r.Cookie("user_id")
	// if token.TokenInput != c.Value {
	// 	WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("token nya bedaa"))
	// 	return
	// }

	// fmt.Println("token input :", token.TokenInput, "dan token cookie :", c.Value)
	tokentod := r.Context().Value("xx").(string)
	if token.TokenInput != tokentod {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("token nya bedaa"))
		return
	}

	dontil := r.Context().Value("id").(entity.Claims)
	if dontil.Role == "admin" {
		eUser, err := u.userService.CheckTokenAdmin(r.Context(), dontil.UserID, token)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		response := map[string]any{
			"user_id": int(eUser.ID),
			"role":    "admin",
			"name":    eUser.ShopName,
			"message": "token benar",
		}

		WriteJSON(w, http.StatusOK, response)
	} else if dontil.Role == "cashier" {
		eUser, err := u.userService.CheckTokenCashier(r.Context(), dontil.UserID, token)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		response := map[string]any{
			"user_id": int(eUser.ID),
			"role":    "cashier",
			"name":    eUser.Username,
			"message": "token benar",
		}

		WriteJSON(w, http.StatusOK, response)
	} else {
		response := map[string]any{
			"message": "apakah kamu hekel",
		}

		WriteJSON(w, http.StatusBadRequest, response)
	}

}
