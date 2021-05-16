package crypto

import "golang.org/x/crypto/bcrypt"

type BCryptComparator struct {
	Cost int
}

func NewComparator(options...int) *BCryptComparator {
	cost := 14
	if len(options) > 0 {
		cost = options[0]
	}
	return &BCryptComparator{Cost: cost}
}

func (b *BCryptComparator) Compare(plaintext []byte, hashed []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashed, plaintext)
	return err == nil, nil
}
func (b *BCryptComparator) Hash(plaintext []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(plaintext, b.Cost)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
