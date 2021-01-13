package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-sso/conf"
	"go-sso/internal/service/apierror"
	"time"
)

var jwtSecret = []byte(conf.GetConfig().Common.AppSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(36 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Issuer:    "qiuplus",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtSecret))
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorMalformed {
				return nil, apierror.ErrTokenMalformed
			} else if ve.Errors == jwt.ValidationErrorExpired {
				// Token is expired
				return nil, apierror.ErrTokenExpired
			} else if ve.Errors == jwt.ValidationErrorNotValidYet {
				return nil, apierror.ErrTokenNotValidYet
			} else {
				return nil, apierror.ErrTokenInvalid
			}
		}
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		fmt.Println(claims, ok)
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, apierror.ErrTokenInvalid
}
