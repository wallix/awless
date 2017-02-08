package logger

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/fatih/color"
)

var DefaultLogger *Logger = &Logger{out: log.New(os.Stdout, "", 0)}

type Logger struct {
	verbose uint32 // atomic
	out     *log.Logger
}

var (
	infoPrefix    = color.GreenString("[info]")
	errorPrefix   = color.RedString("[error]")
	verbosePrefix = color.YellowString("[verbose]")
)

func New(prefix string, flag int) *Logger {
	return &Logger{out: log.New(os.Stdout, prefix, flag)}
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	if l.isVerbose() {
		l.out.Println(prepend(verbosePrefix, fmt.Sprintf(format, v...))...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.out.Println(prepend(infoPrefix, v...)...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.out.Println(prepend(infoPrefix, fmt.Sprintf(format, v...))...)
}

func (l *Logger) Error(v ...interface{}) {
	l.out.Println(prepend(errorPrefix, v...)...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.out.Println(prepend(errorPrefix, fmt.Sprintf(format, v...))...)
}

func (l *Logger) SetVerbose(v bool) {
	if v {
		atomic.StoreUint32(&l.verbose, 1)
	} else {
		atomic.StoreUint32(&l.verbose, 0)
	}
}

func (l *Logger) isVerbose() bool {
	return atomic.LoadUint32(&l.verbose) > 0
}

func Verbosef(format string, v ...interface{}) {
	DefaultLogger.Verbosef(format, v...)
}

func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func prepend(s interface{}, v ...interface{}) []interface{} {
	return append([]interface{}{s}, v...)
}
