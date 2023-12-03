package mux

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	gET    = "GET"
	pOST   = "POST"
	pUT    = "PUT"
	pATCH  = "PATCH"
	dELETE = "DELETE"
)

type GenericHandler interface {
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
type Handler interface {
	GenericHandler
	Search(w http.ResponseWriter, r *http.Request)
}
func RegisterRoutes(r *mux.Router, prefix string, handler GenericHandler, options...bool)  {
	includeDelete := false
	excludePatch := false
	if len(options) > 0 {
		includeDelete = options[0]
	}
	if len(options) > 1 {
		excludePatch = options[1]
	}
	s := r.PathPrefix(prefix).Subrouter()
	s.HandleFunc("/{id}", handler.Load).Methods(gET)
	s.HandleFunc("", handler.Create).Methods(pOST)
	s.HandleFunc("/{id}", handler.Update).Methods(pUT)
	if !excludePatch {
		s.HandleFunc("/{id}", handler.Patch).Methods(pATCH)
	}
	if includeDelete {
		s.HandleFunc("/{id}", handler.Delete).Methods(dELETE)
	}
}

func Register(r *mux.Router, prefix string, handler Handler, options...bool)  {
	includeDelete := false
	includeGetSearch := false
	excludePatch := false
	if len(options) > 0 {
		includeDelete = options[0]
	}
	if len(options) > 1 {
		includeGetSearch = options[1]
	}
	if len(options) > 2 {
		excludePatch = options[2]
	}
	s := r.PathPrefix(prefix).Subrouter()
	if includeGetSearch {
		s.HandleFunc("", handler.Search).Methods(gET)
	}
	s.HandleFunc("/search", handler.Search).Methods(gET)
	s.HandleFunc("/search", handler.Search).Methods(pOST)
	s.HandleFunc("/{id}", handler.Load).Methods(gET)
	s.HandleFunc("", handler.Create).Methods(pOST)
	s.HandleFunc("/{id}", handler.Update).Methods(pUT)
	if !excludePatch {
		s.HandleFunc("/{id}", handler.Patch).Methods(pATCH)
	}
	if includeDelete {
		s.HandleFunc("/{id}", handler.Delete).Methods(dELETE)
	}
}

func HandleFiles(r *mux.Router, prefix string, static string) {
	dir, _ := filepath.Abs(".")
	if _, err := os.Stat(filepath.Join(dir, static)); !os.IsNotExist(err) {
		dir = filepath.Join(dir, static)
	}
	dir = strings.TrimRight(dir, "/")

	h := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
	r.PathPrefix(prefix).Handler(h)
}
