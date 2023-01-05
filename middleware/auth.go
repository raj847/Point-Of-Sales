package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"vandesar/entity"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
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

		tokenString := c.Value

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrInvalidToken
			}
			return []byte("rahasia-perusahaan"), nil
		}

		// Parse methods use this callback function to supply
		// the key for verification.  The function receives the parsed,
		// but unverified Token.  This allows you to use properties in the
		// Header of the token (such as `kid`) to identify which key to use.
		jwtToken, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, keyFunc)
		if err != nil {
			verr, ok := err.(*jwt.ValidationError)
			if ok && errors.Is(verr.Inner, ErrExpiredToken) {
				// return bad request, token expired
				log.Println("token has expired")
			}
      
			// return bad request, invalid token
			log.Println("invalid token")
		}

		claims := jwtToken.Claims.(*entity.Claims)

		fmt.Println("jwt token -> ", tokenString)
		fmt.Println("user id -> ", claims.UserID)

		fmt.Println(tkn.Claims, claims)
		ctx := context.WithValue(r.Context(), "id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
