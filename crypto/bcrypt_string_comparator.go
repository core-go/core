package crypto

import "golang.org/x/crypto/bcrypt"

type BCryptStringComparator struct {
	Cost int
}

func NewStringComparator(options...int) *BCryptStringComparator {
	cost := 14
	if len(options) > 0 {
		cost = options[0]
	}
	return &BCryptStringComparator{Cost: cost}
}

func CompareHashAndPassword(plaintext string, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	return err == nil, nil
}
func Hash(plaintext string) (string, error) {
	return HashWithCost(plaintext, 14)
}
func HashWithCost(plaintext string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func (b *BCryptStringComparator) Compare(plaintext string, hashed string) (bool, error) {
	return CompareHashAndPassword(plaintext, hashed)
}
func (b *BCryptStringComparator) Hash(plaintext string) (string, error) {
	return HashWithCost(plaintext, b.Cost)
}
