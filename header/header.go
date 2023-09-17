package header

import (
	"net/http"
	"time"
)

type Config struct {
	Id string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id"`
	Time string `yaml:"time" mapstructure:"time" json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time"`
	Correlation string `yaml:"correlation" mapstructure:"correlation" json:"correlation,omitempty" gorm:"column:correlation" bson:"correlation,omitempty" dynamodbav:"correlation,omitempty" firestore:"correlation"`
	Constants map[string]string `yaml:"constants" mapstructure:"constants" json:"constants,omitempty" gorm:"column:constants" bson:"constants,omitempty" dynamodbav:"constants,omitempty" firestore:"constants,omitempty"`
}

const format = "2006-01-02T15:04:05.000000-07:00"

type HeaderHandler struct {
	Config Config
	Generate func() string
}
func NewHeaderHandler(config Config, generate func() string) *HeaderHandler {
	return &HeaderHandler{
		Config: config,
		Generate: generate,
	}
}
func (h *HeaderHandler) SetHeader(w http.ResponseWriter, r *http.Request) {
	SetHeader(w, r, h.Config, h.Generate)
}
func (h *HeaderHandler) HandleHeader() func(h http.Handler) http.Handler {
	return func(ht http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			SetHeader(w, r, h.Config, h.Generate)
			ht.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func SetHeader(w http.ResponseWriter, r *http.Request, cfg Config, id func()string) {
	if len(cfg.Correlation) > 0 {
		hdr := r.Header[cfg.Correlation]
		if len(hdr) > 0 {
			w.Header().Set(cfg.Correlation, hdr[0])
		}
	}
	if id != nil && len(cfg.Id) > 0 {
		w.Header().Set(cfg.Id, id())
	}
	if len(cfg.Time) > 0 {
		w.Header().Set(cfg.Time, time.Now().Format(format))
	}
	if len(cfg.Constants) > 0 {
		for k, v := range cfg.Constants {
			w.Header().Set(k, v)
		}
	}
}
