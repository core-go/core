package service

type DefaultUniqueIdGenerator struct {
	ShortId bool
}

func NewUniqueIdGenerator(shortId bool) *DefaultUniqueIdGenerator {
	return &DefaultUniqueIdGenerator{shortId}
}

func (g *DefaultUniqueIdGenerator) Generate() (string, error) {
	if g.ShortId {
		randomId, er3 := ShortId()
		if er3 != nil {
			return "", er3
		}
		return randomId, nil
	} else {
		x := RandomId()
		return x, nil
	}
}
