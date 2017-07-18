package logger

import (
	"io"
	"log"
)

var emergPtr *log.Logger  // Level 0
var alertPtr *log.Logger  // Level 1
var critPtr *log.Logger   // Level 2
var errorPtr *log.Logger  // Level 3
var warnPtr *log.Logger   // Level 4
var noticePtr *log.Logger // Level 5
var infoPtr *log.Logger   // Level 6
var DebugPtr *log.Logger  // Level 7
var TracePtr *log.Logger  // Trace 8

// Create new copy of *log for each Level
func logHandler(h [9]io.Writer) {
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
