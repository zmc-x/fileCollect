package zaplog

import "fileCollect/global"

const (
	INFO int8 = iota
	WARN
	ERROR
	FATAL
)

func GetLogLevel(level int8, msg string) {
	logger := global.Logger
	switch level {
	case INFO:
		logger.Info(msg)
	case WARN:
		logger.Warn(msg)
	case ERROR:
		logger.Error(msg)
	default:
		logger.Fatal(msg)
	}
}