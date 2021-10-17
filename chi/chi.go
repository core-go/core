package chi

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
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
func RegisterRoutes(r *chi.Mux, prefix string, handler GenericHandler, options...bool)  {
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	includeDelete := false
	excludePatch := false
	if len(options) > 0 {
		includeDelete = options[0]
	}
	if len(options) > 1 {
		excludePatch = options[1]
	}
	r.Get(prefix + "{id}", handler.Load)
	r.Post(prefix, handler.Create)
	r.Put(prefix + "{id}", handler.Update)
	if !excludePatch {
		r.Patch(prefix + "{id}", handler.Patch)
	}
	if includeDelete {
		r.Delete(prefix + "{id}", handler.Delete)
	}
}

func Register(r *chi.Mux, prefix string, handler Handler, options...bool)  {
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
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
	if includeGetSearch {
		r.Get(prefix, handler.Search)
	}
	r.Get(prefix + "search", handler.Search)
	r.Post(prefix + "search", handler.Search)
	r.Get(prefix + "{id}", handler.Load)
	r.Post(prefix, handler.Create)
	r.Put(prefix + "{id}", handler.Update)
	if !excludePatch {
		r.Patch(prefix + "{id}", handler.Patch)
	}
	if includeDelete {
		r.Delete(prefix + "{id}", handler.Delete)
	}
}
