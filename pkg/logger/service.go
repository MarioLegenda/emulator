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

var stdOut = true
var logDir = ""

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

func buildBaseLogger(level zapcore.Level, fileName string) *zap.SugaredLogger {
	logFile := fmt.Sprintf("%s/%s", logDir, fileName)

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

func buildInfoLogger() {
	infoLogger = buildBaseLogger(zap.InfoLevel, "info.log")
}

func buildErrorLogger() {
	errorLogger = buildBaseLogger(zap.ErrorLevel, "error.log")
}

func buildWarningLogger() {
	warningLogger = buildBaseLogger(zap.WarnLevel, "warn.log")
}

func BuildLoggers(logDirectory string) {
	logDir = logDirectory
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)

		if err != nil {
			appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create log directory: %s", err.Error()))
		}
	}

	buildInfoLogger()
	buildErrorLogger()
	buildWarningLogger()
}

func Info(msg ...interface{}) {
	if stdOut {
		fmt.Println(msg)
	}

	infoLogger.Info(msg)
}

func Error(msg ...interface{}) {
	if stdOut {
		fmt.Println(msg)
	}

	errorLogger.Error(msg)
}

func Warn(msg ...interface{}) {
	if stdOut {
		fmt.Println(msg)
	}

	warningLogger.Warn(msg)
}
