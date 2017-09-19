# Logger: custom logger package for Go.
Logger package help you easier write log to io.Writer. All methods fully support standard log package, but have extra features.

## Install
For install logger, just run `go get -u github.com/neliseev/logger` or `glide get github.com/neliseev/logger`

### Usage
Simple in main file you should initialise logger.

```golang
package main

import "github.com/neliseev/logger"

var log *logger.Log // Using log subsystem

func init() {
	// Initialization log system
	var err error

	if log, err = logger.NewFileLogger("/var/log/example/main.log", 8); err != nil {
		panic(err)
	}
}

func main() {
    log.Info("Starting main packge")
...
}
```
Logger automatically create all subfolder by defined path.

For re-use current initialised logger, just define var log *logger.Log, example:
```golang
package example

import "github.com/neliseev/logger"

var log *logger.Log // Using log subsystem

func Some() {
    log.Debug("Hello, Wold!")
}
```

#### License
Copyright (C) 2017 Nikita Eliseev

The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software
and associated documentation files (the "Software"), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge, publish, distribute,
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial
portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH
THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
