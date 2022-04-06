package cmd

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(verbose bool) *zap.SugaredLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = localTimeEncoder
	if !verbose {
		encoderConfig.CallerKey = ""
		encoderConfig.LevelKey = ""
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderConfig
	cfg.OutputPaths = []string{"stderr"} // stdoutに出力するとevalでの表示対象となってしまうため

	rowLogger, _ := cfg.Build()

	//nolint:errcheck
	defer rowLogger.Sync() // flushes buffer, if any

	logger := rowLogger.Sugar()
	return logger
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	local, _ := time.LoadLocation("Local")
	enc.AppendString(t.In(local).Format(time.RFC3339))
}
