package utils

import "golang.org/x/crypto/bcrypt"

func Crypt(password string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func Compare(hash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}
