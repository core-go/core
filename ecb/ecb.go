package ecb

import "crypto/cipher"

func NewECBBlockMode(block cipher.Block) cipher.BlockMode {
	return &ECBBlockMode{block}
}
type ECBBlockMode struct {
	block cipher.Block
}
func (b *ECBBlockMode) BlockSize() int {
	return b.block.BlockSize()
}

func (b *ECBBlockMode) CryptBlocks(dst, src []byte) {
	if len(src)%b.block.BlockSize() != 0 {
		panic("crypto/cipher: not full block")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output is smaller than input")
	}
	for len(src) > 0 {
		b.block.Encrypt(dst, src[:b.block.BlockSize()])
		src = src[b.block.BlockSize():]
		dst = dst[b.block.BlockSize():]
	}
}
