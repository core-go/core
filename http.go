package service

import (
	"log"
	"net/http"
)
const InternalServerError = "Internal Server Error"
func Write(w http.ResponseWriter, code int, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if result == nil {
		w.Write([]byte("null"))
		return
	}
	response, err := Marshal(result)
	if err != nil {
		log.Println("cannot marshal of result: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		w.Write(response)
	}
}
func WriteString(w http.ResponseWriter, code int, v string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(v))
}
