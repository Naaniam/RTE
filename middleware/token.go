package middleware

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	JwtExpHour = 0
	JwtExpMin  = 10
	JwtExpSec  = 30
)

// MemberToken fubnction is used to generate a new token for the members
func MemberToken(email string, id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"member": "member",
		"id":     id,
		"email":  email,
		"exp":    time.Now().Local().Add(time.Hour*time.Duration(JwtExpHour) + time.Minute*time.Duration(JwtExpMin) + time.Second*time.Duration(JwtExpSec)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	sToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return sToken, nil
}

func AdminToken(email string, id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": id, "email": email, "admin": "admin", "exp": time.Now().Add(time.Hour * 24 * 30).Unix()})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
