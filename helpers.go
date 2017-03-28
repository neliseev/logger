package logger

import (
	"io"
	"io/ioutil"
	"os"
)

// Open log file and return pointer to FH
func logOpenFile(logFile string) (fh *os.File, err error) {
	defer fh.Close()
	fh, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	return fh, nil
}

// Discard logging destination for reducing log level
func logDiscard(level int, dest [9]io.Writer) [9]io.Writer {
	for ; level < len(dest); level++ {
		dest[level] = ioutil.Discard
	}

	return dest
}
