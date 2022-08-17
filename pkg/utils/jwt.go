package utils

import (
	"banking/pkg/errors"
	"banking/pkg/types"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecret = []byte("miniBankingApp")
var jwtExpiresInMinutes time.Duration = 15

func GenerateAccessToken(login string) (token string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * jwtExpiresInMinutes)),
			IssuedAt:  jwt.NewNumericDate(time.Now())},
	})
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		fmt.Println("Cannot sign access token. Error: ", err)
		return "", err
	}
	return accessTokenString, nil
}

func ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*types.Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", errors.InvalidAccessTokenError
}
