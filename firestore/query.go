package firestore

type Query struct {
	Path     string
	Operator string
	Value    interface{}
}
