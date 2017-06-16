//logger is a wrapper for the log.Logger
//used for
package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	_ logLevel = iota
	LError
	LInfo
	LDebug
)

var (
	l Roga

	Level logLevel
)

type logLevel int

type Roga struct {
	Error *log.Logger
	Info  *log.Logger
	Debug *log.Logger
	Fatal *log.Logger
}

func Load(level logLevel) {
	Level = level
	l.Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	l.Debug = log.New(ioutil.Discard, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Fatal = log.New(os.Stdout, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(v ...interface{}) {
	l.Error.Output(2, fmt.Sprint(v...))
}

func Errorf(f string, v ...interface{}) {
	l.Error.Output(2, fmt.Sprintf(f, v...))
}

func Info(v ...interface{}) {
	l.Info.Output(2, fmt.Sprint(v...))
}

func Infof(f string, v ...interface{}) {
	l.Info.Output(2, fmt.Sprintf(f, v...))
}

func Debug(v ...interface{}) {
	l.Debug.Output(2, fmt.Sprint(v...))
}

func Debugf(f string, v ...interface{}) {
	l.Debug.Output(2, fmt.Sprintf(f, v...))
}

func Fatal(v ...interface{}) {
	l.Fatal.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(f string, v ...interface{}) {
	l.Fatal.Output(2, fmt.Sprintf(f, v...))
	os.Exit(1)
}
