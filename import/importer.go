package impt

import (
	"context"
	"fmt"
	"io"
	"reflect"
)

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

type ErrorHandler struct {
	HandleError func(ctx context.Context, format string, fields map[string]interface{})
	FileName    string
	LineNumber  string
	Map         *map[string]interface{}
}

func NewErrorHandler(logger func(ctx context.Context, format string, fields map[string]interface{}), fileName string, lineNumber string, mp *map[string]interface{}) *ErrorHandler {
	if len(fileName) <= 0 {
		fileName = "filename"
	}
	if len(lineNumber) <= 0 {
		lineNumber = "lineNumber"
	}
	return &ErrorHandler{
		HandleError: logger,
		FileName:    fileName,
		LineNumber:  lineNumber,
		Map:         mp,
	}
}

func (e *ErrorHandler) HandlerError(ctx context.Context, rs interface{}, err []ErrorMessage, i int, fileName string) {
	var ext = make(map[string]interface{})
	if e.Map != nil {
		ext = *e.Map
	}
	if len(e.FileName) > 0 && len(e.LineNumber) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %+v . Error: %+v", rs, err), ext)
	} else if len(e.FileName) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %+v . Error: %+v line: %d", rs, err, i), ext)
	} else if len(e.LineNumber) > 0 {
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %+v . Error: %+v filename:%s", rs, err, fileName), ext)
	} else {
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %+v . Error: %+v filename:%s line: %d", rs, err, fileName, i), ext)
	}
}

func (e *ErrorHandler) HandlerException(ctx context.Context, rs interface{}, err error, i int, fileName string) {
	var ext = make(map[string]interface{})
	if e.Map != nil {
		ext = *e.Map
	}
	if len(e.FileName) > 0 && len(e.LineNumber) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %+v . Error: %+v", rs, err), ext)
	} else if len(e.FileName) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %+v . Error: %+v line: %d", rs, err, i), ext)
	} else if len(e.LineNumber) > 0 {
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %+v . Error: %+v filename:%s", rs, err, fileName), ext)
	} else {
		e.HandleError(ctx, fmt.Sprintf("Error to write: %+v . Error: %v filename: %s line: %d", rs, err, fileName, i), ext)
	}
}

func NewImportRepository(modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error, numLine int) error) error,
	handleException func(ctx context.Context, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	logError func(ctx context.Context, rs interface{}, err []ErrorMessage, i int, fileName string),
	opt ...string,
) *Importer {
	return NewImporter(modelType, transform, write, read, handleException, validate, logError, opt...)
}
func NewImportAdapter(modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error, numLine int) error) error,
	handleException func(ctx context.Context, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	logError func(ctx context.Context, rs interface{}, err []ErrorMessage, i int, fileName string),
	opt ...string,
) *Importer {
	return NewImporter(modelType, transform, write, read, handleException, validate, logError, opt...)
}
func NewImportService(modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error, numLine int) error) error,
	handleException func(ctx context.Context, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	logError func(ctx context.Context, rs interface{}, err []ErrorMessage, i int, fileName string),
	opt ...string,
) *Importer {
	return NewImporter(modelType, transform, write, read, handleException, validate, logError, opt...)
}
func NewImporter(modelType reflect.Type,
	transform func(ctx context.Context, lines []string) (interface{}, error),
	write func(ctx context.Context, data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error, numLine int) error) error,
	handleException func(ctx context.Context, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	handleError func(ctx context.Context, rs interface{}, err []ErrorMessage, i int, fileName string),
	opt ...string,
) *Importer {
	filename := ""
	if len(opt) > 0 {
		filename = opt[0]
	}
	return &Importer{modelType: modelType, Transform: transform, Write: write, Read: read, Validate: validate, HandleError: handleError, HandleException: handleException, Filename: filename}
}

type Importer struct {
	modelType       reflect.Type
	Transform       func(ctx context.Context, lines []string) (interface{}, error)
	Read            func(next func(lines []string, err error, numLine int) error) error
	Write           func(ctx context.Context, data interface{}, endLineFlag bool) error
	Validate        func(ctx context.Context, model interface{}) ([]ErrorMessage, error)
	HandleError     func(ctx context.Context, rs interface{}, err []ErrorMessage, i int, fileName string)
	HandleException func(ctx context.Context, rs interface{}, err error, i int, fileName string)
	Filename        string
}

func (s *Importer) Import(ctx context.Context) (total int, success int, err error) {
	err = s.Read(func(lines []string, err error, numLine int) error {
		if err == io.EOF {
			err = s.Write(ctx, nil, true)
			return nil
		}
		total++
		itemStruct, err := s.Transform(ctx, lines)
		if err != nil {
			return err
		}
		if s.Validate != nil {
			errs, err := s.Validate(ctx, itemStruct)
			if err != nil {
				return err
			}
			if len(errs) > 0 {
				s.HandleError(ctx, itemStruct, errs, numLine, s.Filename)
				return nil
			}
		}
		err = s.Write(ctx, itemStruct, false)
		if err != nil {
			if s.HandleException != nil {
				s.HandleException(ctx, itemStruct, err, numLine, s.Filename)
				return nil
			} else {
				return err
			}
		}
		success++
		return nil
	})
	if err != nil && err != io.EOF {
		return total, success, err
	}
	return total, success, nil
}
