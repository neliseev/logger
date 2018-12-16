package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

// log level constants
const (
	EMERGENCY = iota // 0
	ALERT            // 1
	CRITICAL         // 2
	ERROR            // 3
	WARNING          // 4
	NOTICE           // 5
	INFO             // 6
	DEBUG            // 7
	TRACE            // 8
)

// Log struct represents application logger
type Log struct {
	init  bool
	level int
	fd    *os.File
	d     map[int]*log.Logger // holds writer for each logging level
}

// New constructor for application logger
//
// Fields:
//   path  - path to log file
//   level - log level, possible levels:
//     0 - Emergency, with panic exit.
//     1 - Alert, with exit 2.
//     2 - Critical, with ext 1.
//     3 - Errors.
//     4 - Warnings.
//     5 - Notice.
//     6 - Info.
//     7 - Debug.
//     8 - Trace.
//
// func main() {
//   // initialization application log system
//   var err error
//   logger := logger.Log{}
//   if logger, err = logger.NewLogger("/some/path", 0); err != nil {
//     panic(err)
//   }
// }
func New(cfg Configurer) (*Log, error) {
	level := cfg.LogLevelValue()
	if level > 8 { // validate func parameter
		return nil, fmt.Errorf("incorrect log level, should be from 0 to 8, got: %v", level)
	}

	var (
		logDestination    io.Writer
		logDestinationErr io.Writer
		fd                *os.File
	)

	filePath := cfg.LogFileValue()
	if filePath == "" { // default destinations for logging; write to std out/err
		logDestination = os.Stdout
		logDestinationErr = os.Stderr
	} else { // else try to open file to write
		var err error
		// check if filePath exists
		dirPath := path.Dir(filePath)
		_, err = os.Stat(dirPath)
		if os.IsNotExist(err) { // try to create if directory not exists
			err = os.Mkdir(dirPath, 0750)
		}
		if err == nil {
			fd, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		}
		if err != nil {
			return nil, err
		}
		logDestination = fd
		logDestinationErr = fd
	}

	logger := &Log{
		init:  true,
		fd:    fd,
		level: level,
		d: map[int]*log.Logger{
			EMERGENCY: log.New(logDestinationErr, "Emergency: ", log.Ldate|log.Ltime),
			ALERT:     log.New(logDestinationErr, "Alert: ", log.Ldate|log.Ltime),
			CRITICAL:  log.New(logDestinationErr, "Critical: ", log.Ldate|log.Ltime),
			ERROR:     log.New(logDestinationErr, "Error: ", log.Ldate|log.Ltime),
			WARNING:   log.New(logDestination, "Warning: ", log.Ldate|log.Ltime),
			NOTICE:    log.New(logDestination, "Notice: ", log.Ldate|log.Ltime),
			INFO:      log.New(logDestination, "Info: ", log.Ldate|log.Ltime),
			DEBUG:     log.New(logDestination, "Debug: ", log.Ldate|log.Ltime),
			TRACE:     log.New(logDestination, "Trace: ", log.Ldate|log.Ltime),
		},
	}

	return logger, nil
}

// Close the log file if used
func (log *Log) Close() {
	if log.fd == nil {
		return
	}

	if warn := log.fd.Close(); warn != nil {
		log.println(WARNING, warn)
	}
}

// Emergency wraps println
func (log *Log) Emergency(a ...interface{}) {
	log.println(EMERGENCY, a...)

	panic(fmt.Sprintln(a...))
}

// Alert wraps println
func (log *Log) Alert(a ...interface{}) {
	log.println(ALERT, a...)

	os.Exit(2)
}

// Critical wraps println
func (log *Log) Critical(a ...interface{}) {
	log.println(CRITICAL, a...)

	os.Exit(1)
}

// Error wraps println
func (log *Log) Error(a ...interface{}) {
	log.println(ERROR, a...)
}

// Warning wraps println
func (log *Log) Warning(a ...interface{}) {
	log.println(WARNING, a...)
}

// Notice wraps println
func (log *Log) Notice(a ...interface{}) {
	log.println(NOTICE, a...)
}

// Info wraps println
func (log *Log) Info(a ...interface{}) {
	log.println(INFO, a...)
}

// Debug wraps println
func (log *Log) Debug(a ...interface{}) {
	log.println(DEBUG, a...)
}

// println a message using specified logging level
func (log *Log) println(level int, a ...interface{}) {
	if !log.init { // avoid usage without initialization
		return
	}
	if level < EMERGENCY || level > TRACE {
		level = INFO // use INFO as default fallback
	}
	if level > log.level { // suppress event
		return
	}
	log.d[level].Println(a...) // write message
}

// printlnOnError is a function wrap helper;
// it takes a function that can return an error and additional error description
// if wrapped function return an error we write a message describes event
// and handle behavior accodring to specified logging level further
func (log *Log) printlnOnError(level int, fn func() error, a ...interface{}) {
	err := fn()
	if err == nil { // all is ok
		return
	}

	a = append(a, err)
	log.println(level, a)
}

// WarningOnError wraps printlnOnErr
func (log *Log) WarningOnError(fn func() error, a ...interface{}) {
	log.printlnOnError(WARNING, fn, a...)
}

// ErrorOnError wraps printlnOnErr
func (log *Log) ErrorOnError(fn func() error, a ...interface{}) {
	log.printlnOnError(ERROR, fn, a...)
}

// Trace logging a performance each function (simple profiler), usage defer trace("Message")()
func (log *Log) Trace(msg string) func() {
	if log.level < TRACE { // supress event
		return func() {}
	}

	log.println(TRACE, fmt.Sprintf("Start: %s", msg))
	start := time.Now()

	return func() {
		log.println(TRACE, fmt.Sprintf("Stop: %s, duration: %s", msg, time.Since(start)))
	}
}

// GetPtr returns raw *log.Logger object ptr for specified logging level
func (log *Log) GetPtr(level int) (*log.Logger, error) {
	if !log.init { // avoid usage without initialization
		return nil, errors.New("logger isn't initialized yet")
	}
	if level < EMERGENCY || level > TRACE {
		level = INFO // use INFO as default fallback
	}
	if level > log.level { // supress event
		return nil, fmt.Errorf("current loglevel: %d, got: %d", log.level, level)
	}
	return log.d[level], nil
}

// Configurer represents config object interface
type Configurer interface {
	LogFileValue() string
	LogLevelValue() int
}
