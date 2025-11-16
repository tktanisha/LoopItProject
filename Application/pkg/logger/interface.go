package logger

//go:generate mockgen -source=interface.go -destination=../../internal/mock/mock_logger.go -package=mock

type LoggerInterface interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
	Close()
}
