package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

// log level constants
const (
	EMERGENCY int = iota // 0
	ALERT                // 1
	CRITICAL             // 2
	ERROR                // 3
	WARNING              // 4
	NOTICE               // 5
	INFO                 // 6
	DEBUG                // 7
)

var levels = map[int]string{
	0: "EMERGENCY",
	1: "ALERT",
	2: "CRITICAL",
	3: "ERROR",
	4: "WARNING",
	5: "NOTICE",
	6: "INFO",
	7: "DEBUG",
}

const logEntriesChanSize = 5000

type Logger interface {
	Emergency(a ...interface{})
	Emergencyf(format string, a ...interface{})
	Alert(a ...interface{})
	Alertf(format string, a ...interface{})
	Critical(a ...interface{})
	Criticalf(format string, a ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Warning(a ...interface{})
	Warningf(format string, a ...interface{})
	Notice(a ...interface{})
	Noticef(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Close()
}

// Log struct represents application logger
type Log struct {
	cfg Configurer
	fd  *os.File
	d   map[int]*log.Logger // holds writer for each logging level

	quit       chan struct{}
	entries    chan *jsonLogEntry
	waitGroup  sync.WaitGroup
	httpClient httpClient
}

type jsonLogEntry struct {
	Ts   time.Time `json:"ts"`
	Line string    `json:"line"`
}

type promtailStream struct {
	Labels  string          `json:"labels"`
	Entries []*jsonLogEntry `json:"entries"`
}

type promtailMsg struct {
	Streams []promtailStream `json:"streams"`
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
func New(cfg Configurer) (Logger, error) {
	if cfg.LogLevel() > 8 { // validate func parameter
		return nil, fmt.Errorf("incorrect log level, should be from 0 to 8, got: %v", cfg.LogLevel())
	}

	var (
		err error
		fd  *os.File
	)

	msgOutput := os.Stdout
	errOutput := os.Stderr

	if cfg.LogFile() != "" {
		// check if filePath exists
		dir := path.Dir(cfg.LogFile())
		_, err = os.Stat(dir)
		if os.IsNotExist(err) {
			// try to create if directory not exists
			err = os.Mkdir(dir, 0750)
		} else if err != nil {
			return nil, err
		}

		fd, err = os.OpenFile(cfg.LogFile(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			return nil, err
		}
		msgOutput = fd
		errOutput = fd
	}

	logger := &Log{
		cfg: cfg,
		fd:  fd,
		d: map[int]*log.Logger{
			EMERGENCY: log.New(errOutput, "Emergency: ", log.Ldate|log.Ltime),
			ALERT:     log.New(errOutput, "Alert: ", log.Ldate|log.Ltime),
			CRITICAL:  log.New(errOutput, "Critical: ", log.Ldate|log.Ltime),
			ERROR:     log.New(errOutput, "Error: ", log.Ldate|log.Ltime),
			WARNING:   log.New(msgOutput, "Warning: ", log.Ldate|log.Ltime),
			NOTICE:    log.New(msgOutput, "Notice: ", log.Ldate|log.Ltime),
			INFO:      log.New(msgOutput, "Info: ", log.Ldate|log.Ltime),
			DEBUG:     log.New(msgOutput, "Debug: ", log.Ldate|log.Ltime),
		},
	}

	if cfg.LokiPushURL() != "" {
		logger.quit = make(chan struct{})
		logger.entries = make(chan *jsonLogEntry, logEntriesChanSize)
		logger.httpClient = httpClient{}

		logger.waitGroup.Add(1)
		go logger.run()

	}

	return logger, nil
}

// http.Client wrapper for adding new methods, particularly sendJsonReq
type httpClient struct {
	parent http.Client
}

// A bit more convenient method for sending requests to the HTTP server
func (client *httpClient) sendJsonReq(method, url string, ctype string, reqBody []byte) (resp *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", ctype)
	resp, err = client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBody, nil
}

func (log *Log) run() {
	var batch []*jsonLogEntry
	batchSize := 0
	maxWait := time.NewTimer(log.cfg.LokiBatchWait())

	defer func() {
		if batchSize > 0 {
			log.send(batch)
		}

		log.waitGroup.Done()
	}()

	for {
		select {
		case <-log.quit:
			return
		case entry := <-log.entries:
			batch = append(batch, entry)
			batchSize++
			if batchSize >= 1024 {
				log.send(batch)
				batch = []*jsonLogEntry{}
				batchSize = 0
				maxWait.Reset(log.cfg.LokiBatchWait())
			}
		case <-maxWait.C:
			if batchSize > 0 {
				log.send(batch)
				batch = []*jsonLogEntry{}
				batchSize = 0
			}
			maxWait.Reset(log.cfg.LokiBatchWait())
		}
	}
}

func (log *Log) send(entries []*jsonLogEntry) {
	var streams []promtailStream
	streams = append(streams, promtailStream{
		Labels:  log.cfg.LokiLabels(),
		Entries: entries,
	})

	msg := promtailMsg{Streams: streams}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("promtail.ClientJson: unable to marshal a JSON document: %s\n", err)

		return
	}

	resp, body, err := log.httpClient.sendJsonReq("POST", log.cfg.LokiPushURL(), "application/json", jsonMsg)
	if err != nil {
		fmt.Printf("promtail.ClientJson: unable to send an HTTP request: %s\n", err)

		return
	}

	if resp.StatusCode != 204 {
		fmt.Printf("promtail.ClientJson: Unexpected HTTP status code: %d, response body: %s, request body: %s\n", resp.StatusCode, body, jsonMsg)

		return
	}
}

// Emergency wraps println
func (log *Log) Emergency(a ...interface{}) {
	log.println(EMERGENCY, a...)

	panic(fmt.Sprintln(a...))
}

// Emergencyf wraps formatted Error
func (log *Log) Emergencyf(format string, a ...interface{}) {
	log.printf(EMERGENCY, format, a...)
}

// Alert wraps println
func (log *Log) Alert(a ...interface{}) {
	log.println(ALERT, a...)

	os.Exit(2)
}

// Alertf wraps formatted Error
func (log *Log) Alertf(format string, a ...interface{}) {
	log.printf(ALERT, format, a...)

	os.Exit(2)
}

// Critical wraps println
func (log *Log) Critical(a ...interface{}) {
	log.println(CRITICAL, a...)

	os.Exit(1)
}

// Criticalf wraps formatted Error
func (log *Log) Criticalf(format string, a ...interface{}) {
	log.printf(CRITICAL, format, a...)

	os.Exit(1)
}

// Error wraps println
func (log *Log) Error(a ...interface{}) {
	log.println(ERROR, a...)
}

// Errorf wraps formatted Error
func (log *Log) Errorf(format string, a ...interface{}) {
	log.printf(ERROR, format, a...)
}

// Warning wraps println
func (log *Log) Warning(a ...interface{}) {
	log.println(WARNING, a...)
}

// Warningf wraps formatted Warning
func (log *Log) Warningf(format string, a ...interface{}) {
	log.printf(WARNING, format, a...)
}

// Notice wraps println
func (log *Log) Notice(a ...interface{}) {
	log.println(NOTICE, a...)
}

// Noticef wraps formatted Notice
func (log *Log) Noticef(format string, a ...interface{}) {
	log.printf(NOTICE, format, a...)
}

// Info wraps println
func (log *Log) Info(a ...interface{}) {
	log.println(INFO, a...)
}

// Infof wraps formatted Notice
func (log *Log) Infof(format string, a ...interface{}) {
	log.printf(INFO, format, a...)
}

// Debug wraps println
func (log *Log) Debug(a ...interface{}) {
	log.println(DEBUG, a...)
}

// Debugf wraps formatted Notice
func (log *Log) Debugf(format string, a ...interface{}) {
	log.printf(DEBUG, format, a...)
}

// Close the log file if used
func (log *Log) Close() {
	if log.fd == nil {
		return
	}

	if err := log.fd.Close(); err != nil {
		log.println(WARNING, err)
	}

	if log.cfg.LokiPushURL() != "" {
		close(log.quit)
		log.waitGroup.Wait()
	}
}

func (log *Log) printf(level int, format string, a ...interface{}) {
	if level > log.cfg.LogLevel() { // suppress event
		return
	}

	log.println(level, fmt.Sprintf(format, a...)) // write message
}

// println a message using specified logging level
func (log *Log) println(level int, a ...interface{}) {
	if level > log.cfg.LogLevel() { // suppress event
		return
	}

	log.d[level].Println(a...) // write message

	if log.cfg.LokiPushURL() != "" && log.cfg.LokiSendLevel() <= level {
		line := fmt.Sprintf("%s: %s", levels[level], a)
		log.entries <- &jsonLogEntry{
			Ts:   time.Now(),
			Line: line,
		}
	}
}

// Configurer represents config object interface
type Configurer interface {
	LogFile() string
	LogLevel() int

	LokiPushURL() string
	LokiLabels() string // "{foo=\"bar\"}"
	LokiBatchWait() time.Duration
	LokiSendLevel() int
}
