package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"vandesar/entity"

	"github.com/dgrijalva/jwt-go"
)

func MustAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerType := r.Header.Get("Content-Type")
		c, err := r.Cookie("user_id")

		if err != nil {
			if headerType == "application/json" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse("error unauthorized user id"))
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		claims := &entity.Claims{}

		tkn, err := jwt.ParseWithClaims(c.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasia-perusahaan"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}
		claims = tkn.Claims.(*entity.Claims)

		if claims.Role != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error unauthorized user id"))
			return
		}

		ctx := context.WithValue(r.Context(), "id", claims.AdminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MustCashier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerType := r.Header.Get("Content-Type")
		c, err := r.Cookie("user_id")

		if err != nil {
			if headerType == "application/json" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse("error unauthorized user id"))
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		claims := &entity.Claims{}

		tkn, err := jwt.ParseWithClaims(c.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasia-perusahaan"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}
		claims = tkn.Claims.(*entity.Claims)

		ctx := context.WithValue(r.Context(), "id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerType := r.Header.Get("Content-Type")
		c, err := r.Cookie("user_id")

		if err != nil {
			if headerType == "application/json" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse("error unauthorized user id"))
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		claims := &entity.Claims{}

		tkn, err := jwt.ParseWithClaims(c.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasia-perusahaan"), nil
		})

		if err != nil {
			// fmt.Println(claims)
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		claims = tkn.Claims.(*entity.Claims)
		ctx := context.WithValue(r.Context(), "id", claims.AdminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// package middleware

// import "net/http"

// func Cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func enableCors(w *http.ResponseWriter, r *http.Request) {
// 		dontol := "http://localhost:3000"
// 		(*w).Header().Set("Access-Control-Allow-Origin", dontol)
// 		next.ServeHTTP(w, r)
// 	})
// }
