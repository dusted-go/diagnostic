package log

type Filter interface {
	CanWrite(e event) bool
}

type NoFilter struct{}

func (f *NoFilter) CanWrite(e event) bool {
	return true
}
