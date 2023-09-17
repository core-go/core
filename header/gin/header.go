package header

import (
	"github.com/gin-gonic/gin"
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
func (h *HeaderHandler) SetHeader(ctx *gin.Context) {
	SetHeader(ctx, h.Config, h.Generate)
}
func (h *HeaderHandler) HandleHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		SetHeader(c, h.Config, h.Generate)
		c.Next()
	}
}

func SetHeader(ctx *gin.Context, cfg Config, id func()string) {
	if len(cfg.Correlation) > 0 {
		hdr := ctx.Request.Header[cfg.Correlation]
		if len(hdr) > 0 {
			ctx.Header(cfg.Correlation, hdr[0])
		}
	}
	if id != nil && len(cfg.Id) > 0 {
		ctx.Header(cfg.Id, id())
	}
	if len(cfg.Time) > 0 {
		ctx.Header(cfg.Time, time.Now().Format(format))
	}
	if len(cfg.Constants) > 0 {
		for k, v := range cfg.Constants {
			ctx.Header(k, v)
		}
	}
}
