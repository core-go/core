package service

type Status int

const (
	StatusNotFound     = Status(0)
	StatusSuccess      = Status(1)
	StatusVersionError = Status(2)
	StatusError        = Status(4)
)
