package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level variable using for checking before launch debug and trace
var Level int = 6

// Params - parameters struct.
//
// Fields:
//   LogFile  - path to log file
//   LogLevel - log level, possible levels:
//     0 - Emergency, with panic exit.
//     1 - Alert, with exit 2.
//     2 - Critical, with ext 1.
//     3 - Errors.
//     4 - Warnings.
//     5 - Notice.
//     6 - Info.
//     7 - Debug.
//     8 - Trace.
type Params struct {
	LogFile  string
	LogLevel int
}

func (p *Params) InitLogger() {
	if err := NewLogger(&Params{LogLevel: Level}); err != nil {
		panic(err)
	}
}

// NewLogger() - creating new logger.
func NewLogger(p *Params) error {
	// Default destinations for logging
	dest := [9]io.Writer{
		os.Stderr, // Emerg
		os.Stderr, // Alert
		os.Stderr, // Crit
		os.Stderr, // Err
		os.Stdout, // Warn
		os.Stdout, // Notice
		os.Stdout, // Info
		os.Stdout, // Debug
		os.Stdout, // Trace
	}

	// Change default destination to file if in config defined
	if fh, err := logOpenFile(p.LogFile); fh != nil {
		for i := range dest {
			dest[i] = fh
		}
		// If log file undefined in config, logger will write all into STDOUT/STDERR
	} else if err != nil && p.LogFile != "" {
		logHandler(dest)

		return err
	}

	// Reduce log level
	switch p.LogLevel {
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

		l := new(Params)
		l.Critf("Incorrect log level in config file, defined: %v, possible 0-8", p.LogLevel)
	}

	// Init log.Logger for each Level
	logHandler(dest)

	return nil
}

// Emerg logging a message using Emerg (0) as log level and call panic(fmt string).
func (p *Params) Emerg(args ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintln(args...))
	panic(s)
}

// Emergf logging a message using Emerg (0) as log level and call panic(fmt string).
func (p *Params) Emergf(format string, args ...interface{}) {
	s := emergPtr.Output(2, fmt.Sprintf(format, args...))
	panic(s)
}

// Alert logging a message using Alert (1) as log level and call os.Exit(1).
func (p *Params) Alert(args ...interface{}) {
	alertPtr.Output(2, fmt.Sprintln(args...))
	os.Exit(2)
}

// Alertf logging a message using Alert (1) as log level and call os.Exit(1).
func (p *Params) Alertf(format string, args ...interface{}) {
	alertPtr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Crit logging a message using Crit (2) as log level and call os.Exit(1).
func (p *Params) Crit(args ...interface{}) {
	critPtr.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

// Critf logging a message using Crit (2) as log level and call os.Exit(1).
func (p *Params) Critf(format string, args ...interface{}) {
	critPtr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Err logging a message using Error (3) as log level.
func (p *Params) Err(args ...interface{}) {
	errorPtr.Output(2, fmt.Sprintln(args...))
}

// Errf logging a message using Error (3) as log level.
func (p *Params) Errf(format string, args ...interface{}) {
	errorPtr.Output(2, fmt.Sprintf(format, args...))
}

// Warn logging a message using Warn (4) as log level.
func (p *Params) Warn(args ...interface{}) {
	warnPtr.Println(args)
}

// Warnf logging a message using Warn (4) as log level.
func (p *Params) Warnf(format string, args ...interface{}) {
	warnPtr.Printf(format, args...)
}

// Notice logging a message using Notice (5) as log level.
func (p *Params) Notice(args ...interface{}) {
	noticePtr.Println(args)
}

// Noticef logging a message using Notice (5) as log level.
func (p *Params) Noticef(format string, args ...interface{}) {
	noticePtr.Printf(format, args...)
}

// Info logging a message using Info (6) as log level.
func (p *Params) Info(args ...interface{}) {
	infoPtr.Println(args)
}

// Infof logging a message using Info (6) as log level.
func (p *Params) Infof(format string, args ...interface{}) {
	infoPtr.Printf(format, args...)
}

// Debug logging a message using DEBUG as log level.
func (p *Params) Debug(args ...interface{}) {
	if Level >= 7 {
		DebugPtr.Output(2, fmt.Sprintln(args...))
	}
}

// Debugf logging a message using DEBUG as log level.
func (p *Params) Debugf(format string, args ...interface{}) {
	if Level >= 7 {
		DebugPtr.Output(2, fmt.Sprintf(format, args...))
	}
}

// Trace logging a performance each function, usage defer trace("Message")()
func (p *Params) Trace(msg string) func() {
	if Level == 8 {
		start := time.Now()
		TracePtr.Output(2, fmt.Sprintf("Start: %s", msg))

		return func() {
			TracePtr.Output(2, fmt.Sprintf("Stop: %s, duration: %s", msg, time.Since(start)))
		}
	}

	return func() {}
}
