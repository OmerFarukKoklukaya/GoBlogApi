package utils

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type CustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("JWT_SECRET_KEY")

func ControlToken(rawToken string) *CustomClaims {
	if rawToken == "" {
		fmt.Println("not logged in")
		return nil
	}
	parsedToken, err := jwt.ParseWithClaims(rawToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return parsedToken.Claims.(*CustomClaims)
}

func GenerateToken(authorId uint) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "BlogTest",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	jwtToken.Claims.(*CustomClaims).ID = strconv.Itoa(int(authorId))
	parsedToken, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
	}
	return parsedToken
}
