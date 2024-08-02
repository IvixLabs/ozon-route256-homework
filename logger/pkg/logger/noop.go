package logger

type noop struct {
}

func (*noop) Infow(msg string, keysAndValues ...interface{}) {

}
func (*noop) Warnw(msg string, keysAndValues ...interface{}) {

}
func (*noop) Errorw(msg string, keysAndValues ...interface{}) {

}
func (*noop) Level() Level {

	return ErrorLevel
}
