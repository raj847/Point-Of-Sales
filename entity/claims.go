package entity

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID  uint
	AdminID uint
	Role    string

	jwt.StandardClaims
}
