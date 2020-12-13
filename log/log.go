package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(name string, writer io.Writer) *zap.SugaredLogger {
	conf := zap.NewDevelopmentEncoderConfig()
	conf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(t.UnixNano() / 1e6)
	}
	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncodeCaller = nil

	var lvl zapcore.Level
	if os.Getenv("DEBUG") != "" {
		lvl = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf),
		zapcore.AddSync(writer),
		lvl,
	)

	logger := zap.New(core, zap.AddCaller()).Named(name)
	zap.RedirectStdLog(logger)

	return logger.Sugar()
}
