package reader

type FileType string

const (
	DelimiterType   FileType = "Delimiter"
	FixedlengthType FileType = "Fixedlength"
)

type Reader interface {
	Read(next func(lines string, err error, numLine int) error) error
}
