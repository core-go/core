package service

type UniqueIdGenerator interface {
	Generate() (string, error)
}
