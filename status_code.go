package service

type StatusCode int

const (
	StatusDataNotFoundError = StatusCode(0)
	StatusSuccess           = StatusCode(1)
	StatusError             = StatusCode(2)
	StatusDataVersionError  = StatusCode(4)
)
