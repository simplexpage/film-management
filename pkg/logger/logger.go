package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	LogKeyLevel     LogKey = "level"
	LogKeyMessage   LogKey = "message"
	LogKeySource    LogKey = "source"
	LogKeyTimestamp LogKey = "@timestamp"
	LogKeyTrace     LogKey = "trace"
)

type (
	LogKey string
)

var logLevels = map[string]zapcore.Level{
	"Debug":  zap.DebugLevel,
	"Info":   zap.InfoLevel,
	"Warn":   zap.WarnLevel,
	"Error":  zap.ErrorLevel,
	"DPanic": zap.DPanicLevel,
	"Panic":  zap.PanicLevel,
	"Fatal":  zap.FatalLevel,
}

type Config struct {
	JSON        bool
	Level       string
	Colored     bool
	Development bool
}

func GetZapLogger(config *Config) *zap.Logger {
	level, ok := logLevels[config.Level]
	if !ok {
		level = zap.ErrorLevel
	}

	output := "console"
	if config.JSON {
		output = "json"
	}

	encoder := zapcore.CapitalLevelEncoder
	if config.Colored {
		encoder = zapcore.CapitalColorLevelEncoder
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       config.Development,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          output,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     string(LogKeyMessage),
			LevelKey:       string(LogKeyLevel),
			TimeKey:        string(LogKeyTimestamp),
			NameKey:        "logger",
			CallerKey:      string(LogKeySource),
			StacktraceKey:  string(LogKeyTrace),
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    nil,
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Panic(err)
	}

	return logger
}
