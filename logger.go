package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var level int

var emergPtr *log.Logger
var alertPtr *log.Logger
var critPtr *log.Logger
var errorPtr *log.Logger
var warnPtr *log.Logger
var noticePtr *log.Logger
var infoPtr *log.Logger
var DebugPtr *log.Logger
var TracePtr *log.Logger

func defaultHandler(h [9]io.Writer) {
	emergPtr = log.New(h[0], "Emergency: ", log.Ldate|log.Ltime|log.Lshortfile)
	alertPtr = log.New(h[1], "Alert: ", log.Ldate|log.Ltime|log.Lshortfile)
	critPtr = log.New(h[2], "Critical: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorPtr = log.New(h[3], "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnPtr = log.New(h[4], "Warning: ", log.Ldate|log.Ltime)
	noticePtr = log.New(h[5], "Notice: ", log.Ldate|log.Ltime)
	infoPtr = log.New(h[6], "Info: ", log.Ldate|log.Ltime)
	DebugPtr = log.New(h[7], "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	TracePtr = log.New(h[8], "Trace: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// logDiscard func - discard logging destination for reducing log level.
func discardLevel(level int, dest [9]io.Writer) [9]io.Writer {
	for ; level < len(dest); level++ {
		dest[level] = ioutil.Discard
	}

	return dest
}

type Log struct{}

// Emerg logging a message using Emerg (0) as log level and call panic(fmt string).
func (*Log) Emerg(args ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintln(args...))
	panic(s)
}

// Emergf logging a message using Emerg (0) as log level and call panic(fmt string).
func (*Log) Emergf(format string, args ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintf(format, args...))
	panic(s)
}

// Alert logging a message using Alert (1) as log level and call os.Exit(1).
func (*Log) Alert(args ...interface{}) {
	alertPtr.Output(2, fmt.Sprintln(args...))
	os.Exit(2)
}

// Alertf logging a message using Alert (1) as log level and call os.Exit(1).
func (*Log) Alertf(format string, args ...interface{}) {
	alertPtr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Crit logging a message using Crit (2) as log level and call os.Exit(1).
func (*Log) Crit(args ...interface{}) {
	critPtr.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

// Critf logging a message using Crit (2) as log level and call os.Exit(1).
func (*Log) Critf(format string, args ...interface{}) {
	critPtr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Err logging a message using Error (3) as log level.
func (*Log) Err(args ...interface{}) {
	errorPtr.Output(2, fmt.Sprintln(args...))
}

// Errf logging a message using Error (3) as log level.
func (*Log) Errf(format string, args ...interface{}) {
	errorPtr.Output(2, fmt.Sprintf(format, args...))
}

// Warn logging a message using Warn (4) as log level.
func (*Log) Warn(args ...interface{}) {
	warnPtr.Println(args)
}

// Warnf logging a message using Warn (4) as log level.
func (*Log) Warnf(format string, args ...interface{}) {
	warnPtr.Printf(format, args...)
}

// Notice logging a message using Notice (5) as log level.
func (*Log) Notice(args ...interface{}) {
	noticePtr.Println(args)
}

// Noticef logging a message using Notice (5) as log level.
func (*Log) Noticef(format string, args ...interface{}) {
	noticePtr.Printf(format, args...)
}

// Info logging a message using Info (6) as log level.
func (*Log) Info(args ...interface{}) {
	infoPtr.Println(args)
}

// Infof logging a message using Info (6) as log level.
func (*Log) Infof(format string, args ...interface{}) {
	infoPtr.Printf(format, args...)
}

// Debug logging a message using DEBUG as log level.
func (l *Log) Debug(args ...interface{}) {
	if level >= 7 {
		DebugPtr.Output(2, fmt.Sprintln(args...))
	}
}

// Debugf logging a message using DEBUG as log level.
func (l *Log) Debugf(format string, args ...interface{}) {
	if level >= 7 {
		DebugPtr.Output(2, fmt.Sprintf(format, args...))
	}
}

// Trace logging a performance each function, usage defer trace("Message")()
func (l *Log) Trace(msg string) func() {
	if level == 8 {
		start := time.Now()
		TracePtr.Output(2, fmt.Sprintf("Start: %s", msg))

		return func() {
			TracePtr.Output(2, fmt.Sprintf("Stop: %s, duration: %s", msg, time.Since(start)))
		}
	}

	return func() {}
}
