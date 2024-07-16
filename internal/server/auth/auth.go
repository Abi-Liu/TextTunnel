package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func CompareHashAndPassword(passwordToCompare, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordToCompare), []byte(password))

	if err != nil {
		return false
	}

	return true
}
