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
func (*Log) Emerg(a ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintln(a...))
	panic(s)
}

func (l *Log) EmergOnErr(fn func() error, a ...interface{}) {
	if err := fn(); err != nil {
		a = append(a, err)
		s := emergPtr.Output(2, fmt.Sprintln(a...))
		panic(s)
	}
}

// Emergf logging a message using Emerg (0) as log level and call panic(fmt string).
func (*Log) Emergf(format string, a ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintf(format, a...))
	panic(s)
}

func (l *Log) EmergfOnErr(fn func() error, format string, a ...interface{}) {
	if err := fn(); err != nil {
		format += ": %v"
		a = append(a, err)
		s := emergPtr.Output(2, fmt.Sprintf(format, a...))
		panic(s)
	}
}

// Alert logging a message using Alert (1) as log level and call os.Exit(1).
func (*Log) Alert(a ...interface{}) {
	alertPtr.Output(2, fmt.Sprintln(a...))
	os.Exit(2)
}

func (l *Log) AlertOnErr(fn func() error, a ...interface{}) {
	if err := fn(); err != nil {
		a = append(a, err)
		alertPtr.Output(2, fmt.Sprintln(a...))
		os.Exit(2)
	}
}

// Alertf logging a message using Alert (1) as log level and call os.Exit(1).
func (*Log) Alertf(format string, a ...interface{}) {
	alertPtr.Output(2, fmt.Sprintf(format, a...))
	os.Exit(2)
}

func (l *Log) AlertfOnErr(fn func() error, format string, a ...interface{}) {
	if err := fn(); err != nil {
		format += ": %v"
		a = append(a, err)
		alertPtr.Output(2, fmt.Sprintf(format, a...))
		os.Exit(2)
	}
}

// Crit logging a message using Crit (2) as log level and call os.Exit(1).
func (*Log) Crit(a ...interface{}) {
	critPtr.Output(2, fmt.Sprintln(a...))
	os.Exit(1)
}

func (l *Log) CritOnErr(fn func() error, a ...interface{}) {
	if err := fn(); err != nil {
		a = append(a, err)
		critPtr.Output(2, fmt.Sprintln(a...))
		os.Exit(1)
	}
}

// Critf logging a message using Crit (2) as log level and call os.Exit(1).
func (*Log) Critf(format string, a ...interface{}) {
	critPtr.Output(2, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func (l *Log) CritfOnErr(fn func() error, format string, a ...interface{}) {
	if err := fn(); err != nil {
		format += ": %v"
		a = append(a, err)
		critPtr.Output(2, fmt.Sprintln(a...))
		os.Exit(1)
	}
}

// Err logging a message using Error (3) as log level.
func (*Log) Err(a ...interface{}) {
	errorPtr.Output(2, fmt.Sprintln(a...))
}

func (l *Log) ErrOnErr(fn func() error, a ...interface{}) {
	if err := fn(); err != nil {
		a = append(a, err)
		errorPtr.Output(2, fmt.Sprintln(a...))
	}
}

// Errf logging a message using Error (3) as log level.
func (*Log) Errf(format string, a ...interface{}) {
	errorPtr.Output(2, fmt.Sprintf(format, a...))
}

func (l *Log) ErrfOnErr(fn func() error, format string, a ...interface{}) {
	if err := fn(); err != nil {
		format += ": %v"
		a = append(a, err)
		errorPtr.Output(2, fmt.Sprintf(format, a...))
	}
}

// Warn logging a message using Warn (4) as log level.
func (*Log) Warn(a ...interface{}) {
	warnPtr.Println(a...)
}

func (l *Log) WarnOnErr(fn func() error, a ...interface{}) {
	if err := fn(); err != nil {
		a = append(a, err)
		warnPtr.Println(a...)
	}
}

// Warnf logging a message using Warn (4) as log level.
func (*Log) Warnf(format string, a ...interface{}) {
	warnPtr.Printf(format, a...)
}

func (l *Log) WarnfOnErr(fn func() error, format string, a ...interface{}) {
	if err := fn(); err != nil {
		format += ": %v"
		a = append(a, err)
		warnPtr.Printf(format, a...)
	}
}

// Notice logging a message using Notice (5) as log level.
func (*Log) Notice(a ...interface{}) {
	noticePtr.Println(a...)
}

// Noticef logging a message using Notice (5) as log level.
func (*Log) Noticef(format string, a ...interface{}) {
	noticePtr.Printf(format, a...)
}

// Info logging a message using Info (6) as log level.
func (*Log) Info(a ...interface{}) {
	infoPtr.Println(a...)
}

// Infof logging a message using Info (6) as log level.
func (*Log) Infof(format string, a ...interface{}) {
	infoPtr.Printf(format, a...)
}

// Debug logging a message using DEBUG as log level.
func (l *Log) Debug(a ...interface{}) {
	if level >= 7 {
		DebugPtr.Output(2, fmt.Sprintln(a...))
	}
}

// Debugf logging a message using DEBUG as log level.
func (l *Log) Debugf(format string, a ...interface{}) {
	if level >= 7 {
		DebugPtr.Output(2, fmt.Sprintf(format, a...))
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
