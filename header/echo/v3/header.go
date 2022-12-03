package header

import (
	"github.com/labstack/echo"
	"time"
)

type Config struct {
	Id string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id"`
	App string `yaml:"app" mapstructure:"app" json:"app,omitempty" gorm:"column:app" bson:"app,omitempty" dynamodbav:"app,omitempty" firestore:"app"`
	Time string `yaml:"time" mapstructure:"time" json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time"`
	Correlation string `yaml:"correlation" mapstructure:"correlation" json:"correlation,omitempty" gorm:"column:correlation" bson:"correlation,omitempty" dynamodbav:"correlation,omitempty" firestore:"correlation"`
}

const format = "2006-01-02T15:04:05.000000-07:00"

type HeaderHandler struct {
	Config Config
	App string
	Generate func() string
}
func NewHeaderHandler(config Config, app string, generate func() string) *HeaderHandler {
	return &HeaderHandler{
		Config: config,
		App: app,
		Generate: generate,
	}
}
func (h *HeaderHandler) SetHeader(ctx echo.Context) {
	SetHeader(ctx, h.Config, h.App, h.Generate)
}

func SetHeader(ctx echo.Context, cfg Config, app string, id func()string) {
	if len(cfg.Correlation) > 0 {
		hdr := ctx.Request().Header[cfg.Correlation]
		if len(hdr) > 0 {
			ctx.Response().Header().Set(cfg.Correlation, hdr[0])
		}
	}
	if id != nil && len(cfg.Id) > 0 {
		ctx.Response().Header().Set(cfg.Id, id())
	}
	ctx.Response().Header().Set(cfg.App, app)
	ctx.Response().Header().Set(cfg.Time, time.Now().Format(format))
}
