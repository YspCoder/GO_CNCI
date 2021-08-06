package base

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"sync"
	"time"
)

const (
	LOG_LEVEL = LT_INFO
)

const (
	LT_INVALID = iota
	LT_DEBUG
	LT_INFO
	LT_WARN
	LT_ERROR
)

type LogFile struct {
	lines int
	path  string
	mutex sync.Mutex
}

var sysMutex sync.Mutex

func writeConsole(color ct.Color, text string) {
	sysMutex.Lock()
	defer sysMutex.Unlock()

	ct.ChangeColor(color, true, ct.Black, false)
	fmt.Printf(text)
	ct.ResetColor()
}

func write(console bool, logType int, format string, args ...interface{}) {
	if logType < LOG_LEVEL {
		return
	}

	logTypeKey := "U"
	color := ct.White
	switch logType {
	case LT_DEBUG:
		logTypeKey = "DEBUG"
		color = ct.White
	case LT_INFO:
		logTypeKey = "INFO"
		color = ct.Green
	case LT_WARN:
		logTypeKey = "WARN"
		color = ct.Yellow
	case LT_ERROR:
		logTypeKey = "ERROR"
		color = ct.Red
	}

	now := time.Now()
	var logFormat = fmt.Sprintf("[%s] [%s] %s\n",
		now.Format("2006-01-02 15:04:05"), logTypeKey, format)
	var text = fmt.Sprintf(logFormat, args...)

	if console {
		writeConsole(color, text)
	}
}

func Error(format string, args ...interface{}) {
	write(true, LT_ERROR, format, args...)
}

func Info(format string, args ...interface{}) {
	write(true, LT_INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	write(true, LT_WARN, format, args...)
}

func Debug(format string, args ...interface{}) {
	write(true, LT_DEBUG, format, args...)
}

func Log(format string, args ...interface{}) {
	write(false, LT_INFO, format, args...)
}
