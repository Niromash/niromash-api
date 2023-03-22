package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"niromash-api/api"
	"niromash-api/utils/environment"
	"time"
)

func GenerateAccessTokenOnly(user api.User) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.GetEmail(),
		"exp":   time.Now().Add(api.JwtAccessTokenLifetime).Unix(),
	})

	// Todo: Is a good practice to sign the token with a secret + the user hashed password ?
	signedAccessToken, err := accessToken.SignedString([]byte(environment.GetJWTSecret()))
	if err != nil {
		return "", nil
	}

	return signedAccessToken, nil
}

func GenerateTokenPair(user api.User) ([2]string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.GetEmail(),
		"exp":   time.Now().Add(api.JwtAccessTokenLifetime).Unix(),
	})

	// Todo: Is a good practice to sign the token with a secret + the user hashed password ?
	signedAccessToken, err := accessToken.SignedString([]byte(environment.GetJWTSecret()))
	if err != nil {
		return [2]string{}, nil
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.GetEmail(),
		"type":  "refresh",
		"exp":   time.Now().Add(api.JwtRefreshTokenLifetime).Unix(),
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(environment.GetJWTSecret()))
	if err != nil {
		return [2]string{}, nil
	}

	return [2]string{signedAccessToken, signedRefreshToken}, nil
}
