package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func StartServer(conf ServerConf, handler http.Handler, options ...*tls.Config) {
	log.Println(ServerInfo(conf))
	srv := CreateServer(conf, handler, options...)
	if conf.Secure && len(conf.Key) > 0 && len(conf.Cert) > 0 {
		err := srv.ListenAndServeTLS(conf.Cert, conf.Key)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	} else {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}
}
func Addr(port *int64) string {
	server := ""
	if port != nil && *port >= 0 {
		server = ":" + strconv.FormatInt(*port, 10)
	}
	return server
}
func ServerInfo(conf ServerConf) string {
	if len(conf.Version) > 0 {
		if conf.Port != nil && *conf.Port >= 0 {
			return "Start service: " + conf.Name + " at port " + strconv.FormatInt(*conf.Port, 10) + " with version " + conf.Version
		} else {
			return "Start service: " + conf.Name + " with version " + conf.Version
		}
	} else {
		if conf.Port != nil && *conf.Port >= 0 {
			return "Start service: " + conf.Name + " at port " + strconv.FormatInt(*conf.Port, 10)
		} else {
			return "Start service: " + conf.Name
		}
	}
}
func CreateServer(conf ServerConf, handler http.Handler, options ...*tls.Config) *http.Server {
	addr := Addr(conf.Port)
	srv := http.Server{
		Addr:      addr,
		Handler:   nil,
		TLSConfig: nil,
	}
	if len(options) > 0 && options[0] != nil {
		srv.TLSConfig = options[0]
	}
	if conf.ReadTimeout != nil {
		srv.ReadTimeout = *conf.ReadTimeout
	}
	if conf.ReadHeaderTimeout != nil {
		srv.ReadHeaderTimeout = *conf.ReadHeaderTimeout
	}
	if conf.WriteTimeout != nil {
		srv.WriteTimeout = *conf.WriteTimeout
	}
	if conf.IdleTimeout != nil  {
		srv.IdleTimeout = *conf.IdleTimeout
	}
	if conf.MaxHeaderBytes != nil && *conf.MaxHeaderBytes > 0 {
		srv.MaxHeaderBytes = *conf.MaxHeaderBytes
	}
	srv.Handler = handler
	return &srv
}
