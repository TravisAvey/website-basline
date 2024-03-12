package logging

import (
	"time"

	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
)

var (
	consoleLog log.Logger
	accessLog  log.Logger
	dataLog    log.Logger
)

type LogLevel uint8

const (
	Trace LogLevel = 1
	Debug LogLevel = 2
	Info  LogLevel = 3
	Warn  LogLevel = 4
	Error LogLevel = 5
	Fatal LogLevel = 6
	Panic LogLevel = 7
	None  LogLevel = 8
)

func Setup() {
	consoleLog = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		},
	}

	accessLog = log.Logger{
		Level: log.InfoLevel,
		Writer: &log.FileWriter{
			Filename:     "logs/access.log",
			LocalTime:    true,
			FileMode:     0600,
			MaxSize:      100 * 1024 * 1024,
			MaxBackups:   7,
			EnsureFolder: true,
		},
	}

	dataLog = log.Logger{
		Level: log.InfoLevel,
		Writer: &log.FileWriter{
			Filename:     "logs/data.log",
			MaxSize:      100 * 1024 * 1024,
			MaxBackups:   7,
			LocalTime:    true,
			FileMode:     0600,
			EnsureFolder: true,
		},
	}

	runner := cron.New(cron.WithLocation(time.Local))
	runner.AddFunc("0 0 * * *", func() { accessLog.Writer.(*log.FileWriter).Rotate() })
	runner.AddFunc("0 0 * * *", func() { dataLog.Writer.(*log.FileWriter).Rotate() })
	go runner.Run()
}

func LogAccess(msg string) {
	accessLog.Info().Msg(msg)
}

func LogData(msg string) {
	dataLog.Info().Msg(msg)
}

func SetAccessLogLevel(level LogLevel) {
	switch level {
	case Trace:
		accessLog.Level = log.TraceLevel
	case Debug:
		accessLog.Level = log.DebugLevel
	case Info:
		accessLog.Level = log.InfoLevel
	case Warn:
		accessLog.Level = log.WarnLevel
	}
}

func SetDataLogLevel(level LogLevel) {
	switch level {
	case Trace:
		dataLog.Level = log.TraceLevel
	case Debug:
		dataLog.Level = log.DebugLevel
	case Info:
		dataLog.Level = log.InfoLevel
	case Warn:
		dataLog.Level = log.WarnLevel
	}
}
