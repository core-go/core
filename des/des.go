package des

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

func Encrypt(key []byte, rs []byte, newBlockMode func(block cipher.Block) cipher.BlockMode) (string, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	bm := newBlockMode(block)
	rt := make([]byte, 8)
	bm.CryptBlocks(rt, rs)
	res := hex.EncodeToString(rt)
	return res, nil
}
