package hex

import "encoding/hex"

func Decode(rq []byte) ([]byte, error) {
	ret := make([]byte, hex.DecodedLen(len(rq)))
	_, err := hex.Decode(ret, rq)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
