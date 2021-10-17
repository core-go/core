package handler

import "net/http"

type ViewHandler interface {
	Load(w http.ResponseWriter, r *http.Request)
}
type ViewSearchHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
}
type Handler interface {
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
type SearchHandler interface {
	Handler
	Search(w http.ResponseWriter, r *http.Request)
}
