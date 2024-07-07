package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func StartServer(cfg ServerConfig, handler http.Handler, options ...*tls.Config) {
	log.Println(ServerInfo(cfg))
	srv := CreateServer(cfg, handler, options...)
	if cfg.Secure && len(cfg.Key) > 0 && len(cfg.Cert) > 0 {
		err := srv.ListenAndServeTLS(cfg.Cert, cfg.Key)
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
func ServerInfo(cfg ServerConfig) string {
	if len(cfg.Version) > 0 {
		if cfg.Port != nil && *cfg.Port >= 0 {
			return "Start service: " + cfg.Name + " at port " + strconv.FormatInt(*cfg.Port, 10) + " with version " + cfg.Version
		} else {
			return "Start service: " + cfg.Name + " with version " + cfg.Version
		}
	} else {
		if cfg.Port != nil && *cfg.Port >= 0 {
			return "Start service: " + cfg.Name + " at port " + strconv.FormatInt(*cfg.Port, 10)
		} else {
			return "Start service: " + cfg.Name
		}
	}
}
func CreateServer(cfg ServerConfig, handler http.Handler, options ...*tls.Config) *http.Server {
	addr := Addr(cfg.Port)
	srv := http.Server{
		Addr:      addr,
		Handler:   nil,
		TLSConfig: nil,
	}
	if len(options) > 0 && options[0] != nil {
		srv.TLSConfig = options[0]
	}
	if cfg.ReadTimeout != nil {
		srv.ReadTimeout = *cfg.ReadTimeout
	}
	if cfg.ReadHeaderTimeout != nil {
		srv.ReadHeaderTimeout = *cfg.ReadHeaderTimeout
	}
	if cfg.WriteTimeout != nil {
		srv.WriteTimeout = *cfg.WriteTimeout
	}
	if cfg.IdleTimeout != nil {
		srv.IdleTimeout = *cfg.IdleTimeout
	}
	if cfg.MaxHeaderBytes != nil && *cfg.MaxHeaderBytes > 0 {
		srv.MaxHeaderBytes = *cfg.MaxHeaderBytes
	}
	srv.Handler = handler
	return &srv
}
