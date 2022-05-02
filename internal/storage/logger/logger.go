package logger

import (
	"log"
	"os"
	"strings"
	"sync"
)

type logLevel int

type myLogger struct {
	sync.Mutex
	l *log.Logger
}

const (
	DEBUG logLevel = iota
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
	NONE
)

var (
	std            = myLogger{l: log.New(os.Stderr, "", log.LstdFlags)}
	level logLevel = WARNING
)

var levelNames = []string{
	"[DEBUG] ",
	"[INFO] ",
	"[NOTICE] ",
	"[WARNING] ",
	"[ERROR] ",
	"[CRITICAL] ",
	"",
}

func SetLevel(levelName string) {
	levelName = strings.ToLower(levelName)
	switch levelName {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "notice":
		level = NOTICE
	case "warning":
		level = WARNING
	case "error":
		level = ERROR
	case "critical":
		level = CRITICAL
	case "none":
		level = NONE
	default:
		std.l.Fatalln("Invalid level value provided, allowed values are: debug, info, notice, warning, error, critical and none")
	}
}

func Fatalln(args ...interface{}) {
	std.l.Fatalln(args...)
}

func Debug(args ...interface{}) {
	if level <= DEBUG && level != NONE {
		println(DEBUG, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if level <= DEBUG && level != NONE {
		printf(DEBUG, format, args...)
	}
}

func Info(args ...interface{}) {
	if level <= INFO && level != NONE {
		println(INFO, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if level <= INFO && level != NONE {
		printf(INFO, format, args...)
	}
}

func Warning(args ...interface{}) {
	if level <= WARNING && level != NONE {
		println(WARNING, args...)
	}
}

func Warningf(format string, args ...interface{}) {
	if level <= WARNING && level != NONE {
		printf(WARNING, format, args...)
	}
}

func Error(args ...interface{}) {
	if level <= ERROR && level != NONE {
		println(ERROR, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if level <= ERROR && level != NONE {
		printf(ERROR, format, args...)
	}
}

func println(level logLevel, args ...interface{}) {
	std.Lock()
	defer std.Unlock()

	std.l.SetPrefix(levelNames[level])
	std.l.Println(args...)
	std.l.SetPrefix(levelNames[NONE])
}

func printf(level logLevel, format string, args ...interface{}) {
	std.Lock()
	defer std.Unlock()

	std.l.SetPrefix(levelNames[level])
	std.l.Printf(format, args...)
	std.l.SetPrefix(levelNames[NONE])
}
