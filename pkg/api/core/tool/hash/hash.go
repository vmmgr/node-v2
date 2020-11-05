package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Generate(data string) string {
	hash := sha256.Sum256([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func Verify(data string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)) == nil
}
