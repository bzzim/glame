package token

import (
	"fmt"
	"time"
)
import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	app    string
	secret string
}

func NewJWT(app, secretKey string) *JWT {
	return &JWT{app: app, secret: secretKey}
}

func (r *JWT) SignToken(expiresIn time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"app": r.app,
		"iat": time.Now().Unix(),
		"exp": expiresIn.Unix(),
	})

	return token.SignedString([]byte(r.secret))
}

func (r *JWT) Validate(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(r.secret), nil
	})
	if err != nil {
		return false
	}
	return true
}
