package service

type IdGenerator interface {
	Generate(model interface{}) (int, error)
}
