//Package logger is a wrapper for the log.Logger
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	_ logLevel = iota
	LError
	LInfo
	LDebug
)

var (
	l Roga
)

type logLevel int

type Config struct {
	Level string `json:"level"`
}

type Roga struct {
	Error *log.Logger
	Info  *log.Logger
	Debug *log.Logger
	Fatal *log.Logger

	Level logLevel
}

func Load(c Config) {
	l.Level = getLevel(c.Level)
	l.Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	l.Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Fatal = log.New(os.Stdout, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(v ...interface{}) {
	l.Error.Output(2, fmt.Sprint(v...))
}

func Errorf(f string, v ...interface{}) {
	l.Error.Output(2, fmt.Sprintf(f, v...))
}

func Info(v ...interface{}) {
	if l.Level >= LInfo {
		l.Info.Output(2, fmt.Sprint(v...))
	}
}

func Infof(f string, v ...interface{}) {
	if l.Level >= LInfo {
		l.Info.Output(2, fmt.Sprintf(f, v...))
	}
}

func Debug(v ...interface{}) {
	if l.Level >= LDebug {
		l.Debug.Output(2, fmt.Sprint(v...))
	}
}

func Debugf(f string, v ...interface{}) {
	if l.Level >= LDebug {
		l.Debug.Output(2, fmt.Sprintf(f, v...))
	}
}

func Fatal(v ...interface{}) {
	l.Fatal.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(f string, v ...interface{}) {
	l.Fatal.Output(2, fmt.Sprintf(f, v...))
	os.Exit(1)
}

func getLevel(level string) logLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LDebug
	case "error":
		return LError
	default:
		return LInfo
	}
}
