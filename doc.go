// Logger - implement log subsystem
// Usage:
//   import logger
//
//   var log *logger.Log
//
//   func main() {
//       log.New()
//       test := "Test"
//       log.Emergf("Some: %v", test)
//   }
//
//   TracePtr and DebugPtr exportable for another loggers, type *log.Logger
package logger
