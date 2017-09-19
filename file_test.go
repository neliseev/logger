package logger

import "testing"

const logFile string = "./testing/main.log"
const logLevel int = 8

func TestNewFileLogger(t *testing.T) {
	log, err := NewFileLogger(logFile, logLevel)
	if err != nil {
		t.Error(err)
	}

	log.Trace("Test Trace")
	log.Debugf("Test Debugf: %v", logLevel)
	log.Debug("Test Debug")
	log.Infof("Test Infof: %v", logLevel)
	log.Info("Test Info")
	log.Noticef("Test Noticef: %v", logLevel)
	log.Notice("Test Notice")
	log.Warnf("Test Warnf: %v", logLevel)
	log.Warn("Test Warn")
	log.Errf("Test Errf: %v", logLevel)
	log.Err("Test Err")
}
