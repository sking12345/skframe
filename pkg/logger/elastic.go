package logger

type ElasticLogger struct {
}

func NewElasticLogger() *ElasticLogger {
	return &ElasticLogger{}
}

func (logger *ElasticLogger) Printf(format string, v ...interface{}) {
	println(format, v)
}
