package core_logger

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context")
	}
	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, err
	}

	timestamp := time.Now().UTC().Format("2006-01T15-04-05.000000")
	logPath := filepath.Join(
		config.Folder,
		timestamp+".log",
	)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	// Пишем и в файл, и в stdout
	wcFile := zapcore.AddSync(file)
	wcStdout := zapcore.AddSync(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, wcFile, zapLvl),
		zapcore.NewCore(encoder, wcStdout, zapLvl),
	)

	logger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: logger,
	}, nil
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
	}
}

func (l *Logger) Close() {
	_ = l.Logger.Sync()
}
