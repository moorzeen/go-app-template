package model

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"userId"`
	Login  string    `json:"login"`
	jwt.RegisteredClaims
}
