package service

import (
	"log"
	"net/http"
	"strconv"
)

func StartServer(handler http.Handler, conf ServerConfig) {
	server := Addr(conf.Port)
	log.Println(ServerInfo(conf))
	if err := http.ListenAndServe(server, handler); err != nil {
		panic(err)
	}
}
func Addr(port *int64) string {
	server := ""
	if port != nil && *port >= 0 {
		server = ":" + strconv.FormatInt(*port, 10)
	}
	return server
}
func ServerInfo(conf ServerConfig) string {
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
