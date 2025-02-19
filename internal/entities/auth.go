package entities

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserId    string
	SessionId string
	Version   string
	jwt.StandardClaims
}
