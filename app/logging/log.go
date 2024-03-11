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
			LocalTime:    false,
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
			MaxSize:      50 * 1024 * 1024,
			MaxBackups:   7,
			LocalTime:    false,
			FileMode:     0600,
			EnsureFolder: true,
		},
	}

	runner := cron.New(cron.WithLocation(time.Local))
	runner.AddFunc("0 0 * * *", func() { accessLog.Writer.(*log.FileWriter).Rotate() })
	runner.AddFunc("0 0 * * *", func() { dataLog.Writer.(*log.FileWriter).Rotate() })
	go runner.Run()
}
