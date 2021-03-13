package log

type Filter interface {
	CanWrite(msg string) bool
}

type NoFilter struct{}

func (f *NoFilter) CanWrite(_ string) bool {
	return true
}
