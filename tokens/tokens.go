package tokens

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

var SECRET = os.Getenv("SECRET")

type MyCustomClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func VerifyJWT(token string) error {
	m := MyCustomClaims{}
	t, err := jwt.ParseWithClaims(token, &m, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		return err
	}

	if t.Valid && m.User == "Twiter" {
		return nil
	}
	return errors.New("Unable to verify token")
}

func GenerateToken() string {
	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"kid":  "user",
		"user": "Twitter",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})

	token, err := signer.SignedString([]byte(SECRET))
	if err != nil {
		log.Fatal("Error signing token", err)
	}
	return token
}
