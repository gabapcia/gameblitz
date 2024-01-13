package zap

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Info(msg string, keysAndValues ...any) {
	logger.Infow(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, "error", err)
	logger.Errorw(msg, keysAndValues...)
}

func Panic(err error, msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, "error", err)
	logger.Panicw(msg, keysAndValues...)
}

func Sync() error {
	return logger.Sync()
}

func Start() {
	if logger != nil {
		return
	}

	core := ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), os.Stdout, zap.InfoLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel), zap.AddCallerSkip(1)).Sugar()
}
