package logger

import (
	"emulator/pkg/appErrors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var infoLogger *zap.SugaredLogger
var errorLogger *zap.SugaredLogger
var warningLogger *zap.SugaredLogger

func wrapLumberjack(level zapcore.Level, fileName string) func(core zapcore.Core) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     10,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		level,
	)

	return func(core2 zapcore.Core) zapcore.Core {
		return core
	}
}

func buildBaseLogger(level zapcore.Level, fileName string, logDirectory string) *zap.SugaredLogger {
	logFile := fmt.Sprintf("%s/%s", logDirectory, fileName)

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFile}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.DisableStacktrace = false
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build(zap.WrapCore(wrapLumberjack(level, logFile)))

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create info logger: %s", err.Error()))
	}

	createdLogger := logger.Sugar()

	err = createdLogger.Sync()

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create info logger: %s", err.Error()))
	}

	return createdLogger
}

func buildInfoLogger(logDirectory string) {
	infoLogger = buildBaseLogger(zap.InfoLevel, "info.log", logDirectory)
}

func buildErrorLogger(logDirectory string) {
	errorLogger = buildBaseLogger(zap.ErrorLevel, "error.log", logDirectory)
}

func buildWarningLogger(logDirectory string) {
	warningLogger = buildBaseLogger(zap.WarnLevel, "warn.log", logDirectory)
}

func BuildLoggers(logDirectory string) {
	if _, err := os.Stat(logDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(logDirectory, os.ModePerm)

		if err != nil {
			appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create log directory: %s", err.Error()))
		}
	}

	buildInfoLogger(logDirectory)
	buildErrorLogger(logDirectory)
	buildWarningLogger(logDirectory)
}

func Info(msg ...interface{}) {
	infoLogger.Info(msg)
}

func Error(msg ...interface{}) {
	errorLogger.Error(msg)
}

func Warn(msg ...interface{}) {
	warningLogger.Warn(msg)
}
