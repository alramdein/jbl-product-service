package model

import "github.com/golang-jwt/jwt"

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	NIK      string `json:"nik"`
	FullName string `json:"full_name"`
	jwt.StandardClaims
}
