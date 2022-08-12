package gin

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func NewResponseWriter(rw gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{Body: bytes.NewBufferString(""), ResponseWriter: rw}
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriter) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
