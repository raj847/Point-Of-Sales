package entity

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID  int
	AdminID int
	Role    string
	jwt.StandardClaims
}
