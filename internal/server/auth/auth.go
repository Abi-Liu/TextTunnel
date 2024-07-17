package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func CompareHashAndPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func GetAuthorizationKey(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	arr := strings.Split(authHeader, "Bearer ")
	if len(arr) < 2 {
		return "", errors.New("Invalid authorization key")
	}
	return arr[1], nil
}
