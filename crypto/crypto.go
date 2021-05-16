package crypto

import "crypto/cipher"

func Encrypt(v []byte, iv []byte, cipherBlock cipher.Block, f func([]byte, int) []byte) []byte {
	encrypter := cipher.NewCBCEncrypter(cipherBlock, iv)
	v2 := f(v, cipherBlock.BlockSize())
	bs := make([]byte, len(v2))
	encrypter.CryptBlocks(bs, v2)
	return bs
}

func Decrypt(v []byte, iv []byte, cipherBlock cipher.Block) []byte {
	decrypter := cipher.NewCBCDecrypter(cipherBlock, iv)
	bs := make([]byte, len(v))
	decrypter.CryptBlocks(bs, v)
	return bs
}
