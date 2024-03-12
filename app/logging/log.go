package logging

import (
	"time"

	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
)

var (
	ConsoleLog log.Logger
	AccessLog  log.Logger
	DataLog    log.Logger
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
	ConsoleLog = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		},
	}

	AccessLog = log.Logger{
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

	DataLog = log.Logger{
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
	runner.AddFunc("0 0 * * *", func() { AccessLog.Writer.(*log.FileWriter).Rotate() })
	runner.AddFunc("0 0 * * *", func() { DataLog.Writer.(*log.FileWriter).Rotate() })
	go runner.Run()
}

func SetAccessLogLevel(level LogLevel) {
	switch level {
	case Trace:
		AccessLog.Level = log.TraceLevel
	case Debug:
		AccessLog.Level = log.DebugLevel
	case Info:
		AccessLog.Level = log.InfoLevel
	case Warn:
		AccessLog.Level = log.WarnLevel
	}
}

func SetDataLogLevel(level LogLevel) {
	switch level {
	case Trace:
		DataLog.Level = log.TraceLevel
	case Debug:
		DataLog.Level = log.DebugLevel
	case Info:
		DataLog.Level = log.InfoLevel
	case Warn:
		DataLog.Level = log.WarnLevel
	}
}
