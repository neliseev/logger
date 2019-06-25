package logger

import (
	"fmt"
	"log"
	"os"
	"path"
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
)

// Log struct represents application logger
type Log struct {
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
//
// func main() {
//   // initialization application log system
//   var err error
//   logger := logger.Log{}
//   if logger, err = logger.NewLogger("/some/path", 0); err != nil {
//     panic(err)
//   }
// }
func New(filePath string, level int) (*Log, error) {
	if level > 7 || level < 3 { // validate func parameter
		return nil, fmt.Errorf("incorrect log level, should be between 3 (ERROR) and 7 (DEBUG), got: %v", level)
	}

	var (
		err error
		fd  *os.File
	)

	msgoutput := os.Stdout
	erroutput := os.Stderr

	if filePath != "" {
		// check if filePath exists
		dir := path.Dir(filePath)
		_, err = os.Stat(dir)
		if os.IsNotExist(err) {
			// try to create if directory not exists
			err = os.Mkdir(dir, 0750)
		} else if err != nil {
			return nil, err
		}

		fd, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			return nil, err
		}
		msgoutput = fd
		erroutput = fd
	}

	logger := &Log{
		fd:    fd,
		level: level,
		d: map[int]*log.Logger{
			EMERGENCY: log.New(erroutput, "Emergency: ", log.Ldate|log.Ltime),
			ALERT:     log.New(erroutput, "Alert: ", log.Ldate|log.Ltime),
			CRITICAL:  log.New(erroutput, "Critical: ", log.Ldate|log.Ltime),
			ERROR:     log.New(erroutput, "Error: ", log.Ldate|log.Ltime),
			WARNING:   log.New(msgoutput, "Warning: ", log.Ldate|log.Ltime),
			NOTICE:    log.New(msgoutput, "Notice: ", log.Ldate|log.Ltime),
			INFO:      log.New(msgoutput, "Info: ", log.Ldate|log.Ltime),
			DEBUG:     log.New(msgoutput, "Debug: ", log.Ldate|log.Ltime),
		},
	}

	return logger, nil
}

// Close the log file if used
func (log *Log) Close() {
	if log.fd == nil {
		return
	}

	if err := log.fd.Close(); err != nil {
		log.println(WARNING, err)
	}
}

// Emergency wraps println
func (log *Log) Emergency(a ...interface{}) {
	log.println(EMERGENCY, a...)

	panic(fmt.Sprintln(a...))
}

// Emergencyf wraps formatted Error
func (log *Log) Emergencyf(format string, a ...interface{}) {
	log.Emergency(fmt.Sprintf(format, a...))
}

// Alert wraps println
func (log *Log) Alert(a ...interface{}) {
	log.println(ALERT, a...)

	os.Exit(2)
}

// Alertf wraps formatted Error
func (log *Log) Alertf(format string, a ...interface{}) {
	log.Alert(fmt.Sprintf(format, a...))
}

// Critical wraps println
func (log *Log) Critical(a ...interface{}) {
	log.println(CRITICAL, a...)

	os.Exit(1)
}

// Criticalf wraps formatted Error
func (log *Log) Criticalf(format string, a ...interface{}) {
	log.Critical(fmt.Sprintf(format, a...))
}

// Error wraps println
func (log *Log) Error(a ...interface{}) {
	log.println(ERROR, a...)
}

// Errorf wraps formatted Error
func (log *Log) Errorf(format string, a ...interface{}) {
	log.Error(fmt.Sprintf(format, a...))
}

// Warning wraps println
func (log *Log) Warning(a ...interface{}) {
	log.println(WARNING, a...)
}

// Warningf wraps formatted Warning
func (log *Log) Warningf(format string, a ...interface{}) {
	log.Warning(fmt.Sprintf(format, a...))
}

// Notice wraps println
func (log *Log) Notice(a ...interface{}) {
	log.println(NOTICE, a...)
}

// Noticef wraps formatted Notice
func (log *Log) Noticef(format string, a ...interface{}) {
	log.Notice(fmt.Sprintf(format, a...))
}

// Info wraps println
func (log *Log) Info(a ...interface{}) {
	log.println(INFO, a...)
}

// Infof wraps formatted Notice
func (log *Log) Infof(format string, a ...interface{}) {
	log.Info(fmt.Sprintf(format, a...))
}

// Debug wraps println
func (log *Log) Debug(a ...interface{}) {
	log.println(DEBUG, a...)
}

// Debugf wraps formatted Notice
func (log *Log) Debugf(format string, a ...interface{}) {
	log.Debug(fmt.Sprintf(format, a...))
}

// println a message using specified logging level
func (log *Log) println(level int, a ...interface{}) {
	if level > log.level { // suppress event
		return
	}

	log.d[level].Println(a...) // write message
}
