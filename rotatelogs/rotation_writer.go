package rotatelogs

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"io"
	"time"
)

func GetWriter(logLocation string, rotationTime time.Duration) (io.Writer, func() error){
	writer, _ := rotatelogs.New(
		logLocation,
		rotatelogs.WithRotationTime(rotationTime),
	)
	return writer, writer.Close
}
