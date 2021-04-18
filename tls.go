package service

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
)

func StartServerTLS(r http.Handler, conf ServerConfig, cert, key []byte) {
	var err error
	server := GetServerAddress(conf.Port)
	log.Println(GetStartMessage(conf))
	s := &http.Server{Addr: server, Handler: r, TLSConfig: &tls.Config{}}
	s.TLSConfig.NextProtos = make([]string, 1)
	if !strSliceContains(s.TLSConfig.NextProtos, "http/1.1") {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "http/1.1")
	}
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)
	s.TLSConfig.Certificates[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}
	ln, err := net.Listen("tcp", server)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = ln.Close()
		if err != nil {
			panic(err)
		}
	}()
	tlsListener := tls.NewListener(ln, s.TLSConfig)
	err = s.Serve(tlsListener)
	if err != nil {
		panic(err)
	}
}

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
