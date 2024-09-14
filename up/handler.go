package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const contentTypeHeader = "Content-Type"

func NewHandler(service UploadService, logError func(context.Context, string, ...map[string]interface{}),
	keyFile string, generate func(ctx context.Context) (string, error), opts ...int,
) *Handler {
	idIndex := 1
	if len(opts) > 0 && opts[0] >= 0 {
		idIndex = opts[0]
	}
	if len(keyFile) == 0 {
		keyFile = "file"
	}
	return &Handler{Service: service, LogError: logError,
		KeyFile: keyFile, generateId: generate, idIndex: idIndex,
	}
}

type Handler struct {
	Service    UploadService
	LogError   func(context.Context, string, ...map[string]interface{})
	KeyFile    string
	generateId func(ctx context.Context) (string, error)
	idIndex    int
}

func (u *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredString(w, r, u.idIndex)
	if len(id) > 0 {
		err := r.ParseMultipartForm(200000)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		formdata := r.MultipartForm // ok, no problem so far, read the Form data

		//get the *fileheaders
		files := formdata.File[u.KeyFile] // grab the filenames
		_, handler, _ := r.FormFile(u.KeyFile)
		contentType := handler.Header.Get(contentTypeHeader)
		if len(contentType) == 0 {
			contentType = getExt(handler.Filename)
		}
		generateStr, err := u.generateId(r.Context())
		if err != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		var list []UploadData
		for i, _ := range files { // loop through the files one by one
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			defer file.Close()
			out := bytes.NewBuffer(nil)

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				if u.LogError != nil {
					u.LogError(r.Context(), err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			bytes := out.Bytes()
			name := generateStr + "_" + files[i].Filename
			list = append(list, UploadData{name, bytes})
		}
		rs, err := u.Service.UploadImage(id, list, contentType, r)
		if err != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, rs)
	}
}

func (u *Handler) UploadGallery(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredString(w, r, u.idIndex)
	if len(id) > 0 {
		err1 := r.ParseMultipartForm(32 << 20)
		if err1 != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		file, handler, err2 := r.FormFile(u.KeyFile)
		if err2 != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		bufferFile := bytes.NewBuffer(nil)
		_, err3 := io.Copy(bufferFile, file)
		if err3 != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err3.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err3.Error(), http.StatusInternalServerError)
			}
			return
		}

		defer file.Close()
		bytes := bufferFile.Bytes()
		contentType := handler.Header.Get(contentTypeHeader)
		if len(contentType) == 0 {
			contentType = getExt(handler.Filename)
		}
		generateStr, _ := u.generateId(r.Context())
		name := generateStr + "_" + handler.Filename
		uploadFile := Upload{Source: r.FormValue("source"), Data: bytes, Name: name, Id: id, Type: contentType}
		rs, err5 := u.Service.UploadGallery(uploadFile, r)

		if err5 != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err5.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err5.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, rs)
	}

}

func (u *Handler) UploadCover(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredString(w, r, u.idIndex)
	if len(id) > 0 {
		err := r.ParseMultipartForm(200000)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		formdata := r.MultipartForm // ok, no problem so far, read the Form data

		//get the *fileheaders
		files := formdata.File[u.KeyFile] // grab the filenames
		_, handler, _ := r.FormFile(u.KeyFile)
		contentType := handler.Header.Get(contentTypeHeader)
		if len(contentType) == 0 {
			contentType = getExt(handler.Filename)
		}
		generateStr, _ := u.generateId(r.Context())
		var list []UploadData
		for i, _ := range files { // loop through the files one by one
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			defer file.Close()
			out := bytes.NewBuffer(nil)

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				if u.LogError != nil {
					u.LogError(r.Context(), err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			bytes := out.Bytes()
			name := generateStr + "_" + files[i].Filename
			list = append(list, UploadData{name, bytes})
		}
		rs, err := u.Service.UploadCover(id, list, contentType, r)
		if err != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, rs)
	}
}

func (u *Handler) DeleteGalleryFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	id := GetRequiredString(w, r, u.idIndex)
	if len(id) > 0 {
		result, err4 := u.Service.DeleteGalleryFile(id, url, r)
		if err4 != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err4.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err4.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, result)
	}
}

func (u *Handler) GetGallery(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredString(w, r, u.idIndex)
	if len(id) > 0 {
		result, err := u.Service.GetGallery(id, r)
		if err != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, result)
	}
}

func (u *Handler) UpdateGallery(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredString(w, r, u.idIndex)
	if len(id) > 0 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := buf.String()
		var body []UploadInfo
		json.NewDecoder(strings.NewReader(s)).Decode(&body)

		result, err4 := u.Service.UpdateGallery(body, id, r)
		if err4 != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err4.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err4.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, result)
	}
}

func respond(w http.ResponseWriter, code int, result interface{}) {
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
func GetRequiredString(w http.ResponseWriter, r *http.Request, options ...int) string {
	p := GetString(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}
func GetString(r *http.Request, options ...int) string {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
	s := r.URL.Path
	params := strings.Split(s, "/")
	i := len(params) - 1 - offset
	if i >= 0 {
		return params[i]
	} else {
		return ""
	}
}
