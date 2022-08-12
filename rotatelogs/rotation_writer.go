package rotatelogs

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"time"
)

func GetWriter(logLocation string, rotationTime time.Duration, maxSize int64) (io.Writer, func() error) {
	opts := make([]rotatelogs.Option, 0)

	if maxSize > 0 {
		opts = append(opts, rotatelogs.WithRotationSize(maxSize))
	}
	opts = append(opts, rotatelogs.WithRotationTime(rotationTime))
	writer, _ := rotatelogs.New(
		logLocation,
		opts...,
	)
	return writer, writer.Close
}
