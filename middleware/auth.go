package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"vandesar/entity"

	"github.com/dgrijalva/jwt-go"
)

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

		fmt.Println(c.Value)

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

		ctx := context.WithValue(r.Context(), "id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
