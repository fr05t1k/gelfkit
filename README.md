# GelfKit

GelfKit allow to use [go-kit logger](https://github.com/go-kit/kit/tree/master/log) with [GELF](https://docs.graylog.org/en/3.0/pages/gelf.html) protocol

# Example
```go
package main

import (
	"fmt"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"github.com/fr05t1k/gelfkit"
)

func main() {
	gelfWriter, _ := gelf.NewWriter("localhost:12201")
	logger, _ := gelfkit.NewGelfLogger(gelfWriter)
	logger = log.With(logger, "caller", log.Caller(4))
	
	logger.Log("msg", "Hello world")
	
}
```
