package logger

type Stub struct {
}

func NewStubLogger() *Stub {
	return &Stub{}
}

func (l *Stub) Info(_ ...any) {
}
func (l *Stub) Warn(_ ...any) {
}

func (l *Stub) Error(_ ...any) {
}
