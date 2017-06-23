//Package logger is a wrapper for the log.Logger
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	_ logLevel = iota
	LDebug
	LInfo
	LError
	LOff
)

var (
	l Roga

	defaultWriter io.Writer
	errorWriter   io.Writer

	logFile *os.File = nil
)

type logLevel int

type Config struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type Roga struct {
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
	Fatal *log.Logger

	Level logLevel
}

func Load(c Config) {
	//open file if there is one
	if file, err := os.OpenFile(c.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		logFile = file
		defaultWriter = io.MultiWriter(os.Stdout, logFile)
		errorWriter = io.MultiWriter(os.Stderr, logFile)
	} else {
		defaultWriter = os.Stdout
		errorWriter = os.Stderr
	}

	//grab the loglevel from config (LInfo is default)
	l.Level = getLevel(c.Level)

	//setup all the loggers
	l.Debug = log.New(defaultWriter, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Info = log.New(defaultWriter, "[INFO] ", log.Ldate|log.Ltime)
	l.Error = log.New(errorWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Fatal = log.New(errorWriter, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(v ...interface{}) {
	if l.Level <= LError {
		l.Error.Output(2, fmt.Sprint(v...))
	}
}

func Errorf(f string, v ...interface{}) {
	if l.Level <= LError {
		l.Error.Output(2, fmt.Sprintf(f, v...))
	}
}

func Info(v ...interface{}) {
	if l.Level <= LInfo {
		l.Info.Output(2, fmt.Sprint(v...))
	}
}

func Infof(f string, v ...interface{}) {
	if l.Level <= LInfo {
		l.Info.Output(2, fmt.Sprintf(f, v...))
	}
}

func Debug(v ...interface{}) {
	if l.Level <= LDebug {
		l.Debug.Output(2, fmt.Sprint(v...))
	}
}

func Debugf(f string, v ...interface{}) {
	if l.Level <= LDebug {
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

func Close() {
	//check if we have a log file
	if logFile != nil {
		l.Info.Output(1, "closing logfile...")
		//attempt to close the logfile
		if err := logFile.Close(); err != nil {
			panic(err)
		}
	}
}

func getLevel(level string) logLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LDebug
	case "error":
		return LError
	case "off":
		return LOff
	default:
		return LInfo
	}
}
