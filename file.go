package logger

import (
	"fmt"
	"io"
	"os"
	"path"
)

// NewFileLogger func - creating new file logger.
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
//
// Example, for reusing current pointer to this logger, just define in package var log *logger.Log:
// var log *logger.Log // Using log subsystem
//
// func init() {
//   // Initialization log system
//   var err error
//
//   if log, err = logger.NewFileLogger(environment.GetLogFile(), environment.GetLogLevel()); err != nil {
//     panic(err)
//   }
//}
func NewFileLogger(logFile string, logLevel int) (*Log, error) {
	if logLevel > 8 {
		return nil, fmt.Errorf("incorrect log level, should be from 0 to 8, got: %v", logLevel)
	}
	level = logLevel

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
	if fh, err := openLogFile(logFile); fh != nil {
		for i := range dest {
			dest[i] = fh
		}
		// If log file undefined in config, logger will write all into STDOUT/STDERR
	} else if err != nil && logFile != "" {
		return nil, err
	}

	// Reduce log level
	switch logLevel {
	case 8:
		dest = discardLevel(9, dest)
	case 7:
		dest = discardLevel(8, dest)
	case 6:
		dest = discardLevel(7, dest)
	case 5:
		dest = discardLevel(6, dest)
	case 4:
		dest = discardLevel(5, dest)
	case 3:
		dest = discardLevel(4, dest)
	case 2:
		dest = discardLevel(3, dest)
	case 1:
		dest = discardLevel(2, dest)
	case 0:
		dest = discardLevel(1, dest)
	}

	// Init log.Logger for each Level
	defaultHandler(dest)

	return new(Log), nil
}

// openLogFile - creating if file not exist and opening it for writes.
func openLogFile(logFile string) (*os.File, error) {
	if err := os.MkdirAll(path.Dir(logFile), 0755); err != nil {
		return nil, err
	}

	fh, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	return fh, nil
}
