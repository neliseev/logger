package logger

import (
	"io"
	"os"
	"fmt"
	"time"
)

// Exporting to core/cfg
type Log struct {
	LogFile   string
	LogLevel  int
}

// Level variable using for checking before launch debug and trace
var Level int = 6

// New() Initializing log io writer
func New(c *Log) error {
	// Default destinations for logging
	dest := [9]io.Writer{
		os.Stderr,  // Emerg
		os.Stderr,  // Alert
		os.Stderr,  // Crit
		os.Stderr,  // Err
		os.Stdout,  // Warn
		os.Stdout,  // Notice
		os.Stdout,  // Info
		os.Stdout,  // Debug
		os.Stdout,  // Trace
	}

	// Change default destination to file if in config defined
	if fh, err := logOpenFile(c.LogFile); fh != nil {
		for i := range dest {
			dest[i] = fh
		}
	// If log file undefined in config, logger will write all into STDOUT/STDERR
	} else if err != nil && c.LogFile != "" {
		logHandler(dest)

		return err
	}

	// Reduce log level
	switch c.LogLevel {
	case 8:
		dest = logDiscard(9, dest)
	case 7:
		dest = logDiscard(8, dest)
	case 6:
		dest = logDiscard(7, dest)
	case 5:
		dest = logDiscard(6, dest)
	case 4:
		dest = logDiscard(5, dest)
	case 3:
		dest = logDiscard(4, dest)
	case 2:
		dest = logDiscard(3, dest)
	case 1:
		dest = logDiscard(2, dest)
	case 0:
		dest = logDiscard(1, dest)
	default:
		logHandler(dest)

		l := new(Log)
		l.Critf("Incoreect log level in config file, defined: %v, possible 0-8", c.LogLevel)
	}

	// Init log.Logger for each Level
	logHandler(dest)

	return nil
}

// Emerg logging a message using Emerg (0) as log level and call panic(fmt string).
func (l *Log) Emerg(args ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintln(args...))
	panic(s)
}

// Emergf logging a message using Emerg (0) as log level and call panic(fmt string).
func (l *Log) Emergf(format string, args ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintf(format, args...))
	panic(s)
}

// Alert logging a message using Alert (1) as log level and call os.Exit(1).
func (l *Log) Alert(args ...interface{}) {
	alertPtr.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

// Alertf logging a message using Alert (1) as log level and call os.Exit(1).
func (l *Log) Alertf(format string, args ...interface{}) {
	alertPtr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Crit logging a message using Crit (2) as log level and call os.Exit(1).
func (l *Log) Crit(args ...interface{}) {
	critPtr.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

// Critf logging a message using Crit (2) as log level and call os.Exit(1).
func (l *Log) Critf(format string, args ...interface{}) {
	critPtr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Err logging a message using Error (3) as log level.
func (l *Log) Err(args ...interface{}) {
	errorPtr.Output(2, fmt.Sprintln(args...))
}

// Errf logging a message using Error (3) as log level.
func (l *Log) Errf(format string, args ...interface{}) {
	errorPtr.Output(2, fmt.Sprintf(format, args...))
}

// Warn logging a message using Warn (4) as log level.
func (l *Log) Warn(args ...interface{}) {
	warnPtr.Println(args)
}

// Warnf logging a message using Warn (4) as log level.
func (l *Log) Warnf(format string, args ...interface{}) {
	warnPtr.Printf(format, args...)
}

// Notice logging a message using Notice (5) as log level.
func (l *Log) Notice(args ...interface{}) {
	noticePtr.Println(args)
}

// Noticef logging a message using Notice (5) as log level.
func (l *Log) Noticef(format string, args ...interface{}) {
	noticePtr.Printf(format, args...)
}

// Info logging a message using Info (6) as log level.
func (l *Log) Info(args ...interface{}) {
	infoPtr.Println(args)
}

// Infof logging a message using Info (6) as log level.
func (l *Log) Infof(format string, args ...interface{}) {
	infoPtr.Printf(format, args...)
}

// 7 - Debug logging a message using DEBUG as log level.
func (l *Log) Debug(args ...interface{}) {
	if Level >= 7 {
		DebugPtr.Output(2, fmt.Sprintln(args...))
	}
}

// 7 - Debugf logging a message using DEBUG as log level.
func (l *Log) Debugf(format string, args ...interface{}) {
	if Level >= 7 {
		DebugPtr.Output(2, fmt.Sprintf(format, args...))
	}
}

// Trace logging a performance each function, usage defer trace("Message")()
func (l *Log) Trace(msg string) func() {
	if Level == 8 {
		start := time.Now()
		TracePtr.Output(2, fmt.Sprintf("Start: %s", msg))

		return func() {
			TracePtr.Output(2, fmt.Sprintf("Stop: %s, duration: %s", msg, time.Since(start)))
		}
	}

	return func() {}
}
