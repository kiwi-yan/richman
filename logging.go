package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LogError LogLevel = iota
	LogWarning
	LogNotice
	LogInfo
	LogDebug
	LogTrace
)

// copy from log/log.go
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type Logger struct {
	logger *log.Logger
	level  LogLevel
}

func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		logger: log.New(out, prefix, flag),
		level:  LogDebug,
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Logger) Output(calldepth int, s string) error {
	return l.logger.Output(calldepth+1, s)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Print(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
}

func (l *Logger) Println(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
}

var std = New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)

func Default() *Logger { return std }

var (
	SetOutput = std.SetOutput
	Println   = std.Println
)

func testLogger() {
	Println("Hello")
	// rand.Seed(3)
	// fmt.Println("My favorite number is", rand.Intn(10))
}
