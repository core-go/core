package service

import "crypto/tls"

func CreateTLSConfig(cert, key []byte) (*tls.Config, error) {
	c := &tls.Config{}
	c.NextProtos = make([]string, 1)
	if !strSliceContains(c.NextProtos, "http/1.1") {
		c.NextProtos = append(c.NextProtos, "http/1.1")
	}
	cer, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	c.Certificates = make([]tls.Certificate, 1)
	c.Certificates[0] = cer
	return c, nil
}
func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
