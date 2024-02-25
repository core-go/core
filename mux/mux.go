package mux

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
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

func HandleWithSecurity(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), check func(next http.Handler) http.Handler, authorize func(next http.Handler, privilege string, action int32) http.Handler, menuId string, action int32, methods ...string) *mux.Route {
	finalHandler := http.HandlerFunc(f)
	funcAuthorize := func(next http.Handler) http.Handler {
		return authorize(next, menuId, action)
	}
	return r.Handle(path, check(funcAuthorize(finalHandler))).Methods(methods...)
}

type SecurityHandler struct {
	Check func(next http.Handler) http.Handler
	Authorize func(next http.Handler, privilege string, action int32) http.Handler
}
func NewSecurityHandler(check func(next http.Handler) http.Handler, authorize func(next http.Handler, privilege string, action int32) http.Handler) *SecurityHandler {
	return &SecurityHandler{Check: check, Authorize: authorize}
}
func (h *SecurityHandler) Handle(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), menuId string, action int32, methods ...string) *mux.Route {
	finalHandler := http.HandlerFunc(f)
	funcAuthorize := func(next http.Handler) http.Handler {
		return h.Authorize(next, menuId, action)
	}
	return r.Handle(path, h.Check(funcAuthorize(finalHandler))).Methods(methods...)
}

type CookieHandler struct {
	Check func(next http.Handler, skipRefreshTTL bool) http.Handler
	Authorize func(next http.Handler, privilege string, action int32) http.Handler
	Prefix string
	Router *mux.Router
}

func NewCookieHandler(r *mux.Router, check func(next http.Handler, skipRefreshTTL bool) http.Handler, authorize func(next http.Handler, privilege string, action int32) http.Handler, opts...string) *CookieHandler {
	ch := &CookieHandler{Router: r, Check: check, Authorize: authorize}
	if len(opts) > 0 && len(opts[0]) > 0 {
		ch.Router = ch.Router.PathPrefix(opts[0]).Subrouter()
		return ch
	}
	return ch
}
func (h *CookieHandler) CheckPrivilege(path string, f http.HandlerFunc, privilege string, action int32, methods ...string) *mux.Route {
	finalHandler := f
	authorize := func(next http.Handler) http.Handler {
		return h.Authorize(next, privilege, action)
	}
	return h.Router.Handle(path, h.Check(authorize(finalHandler), false)).Methods(methods...)
}

func (h *CookieHandler) Handle(path string, f http.HandlerFunc, methods ...string) *mux.Route {
	return h.Router.Handle(path, h.Check(f, true)).Methods(methods...)
}
