package telphin

// FieldLogger interface
type FieldLogger interface {
	Logger
	WithField(string, interface{}) FieldLogger
	WithFields(map[string]interface{}) FieldLogger
}

// Logger interface
// TODO: remove logrus
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
}
