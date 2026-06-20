package helper

import "crypto/sha256"
import "fmt"

// HashPassword - using SHA256 to avoid bcrypt external dependency issue in this env.
// In production, swap with golang.org/x/crypto/bcrypt.
func HashPassword(plain string) string {
	h := sha256.New()
	h.Write([]byte(plain))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CheckPassword(plain, hashed string) bool {
	return HashPassword(plain) == hashed
}
