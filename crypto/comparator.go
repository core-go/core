package crypto

type Comparator interface {
	Compare(plaintext []byte, hashed []byte) (bool, error)
	Hash(plaintext []byte) ([]byte, error)
}
