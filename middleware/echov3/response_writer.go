package echo

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{Body: bytes.NewBufferString(""), ResponseWriter: rw}
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}
