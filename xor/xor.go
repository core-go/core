package xor

import "errors"

func XOR(input1, input2 []byte) ([]byte, error) {
	if len(input1) != len(input2) {
		return nil, errors.New("mismatched lengths")
	}
	ret := make([]byte, len(input1))
	for i := 0; i < len(input1); i++ {
		ret[i] = input1[i] ^ input2[i]
	}
	return ret, nil
}
