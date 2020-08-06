# GelfKit

GelfKit allow to use [go-kit logger](https://github.com/go-kit/kit/tree/master/log) with [GELF](https://docs.graylog.org/en/3.0/pages/gelf.html) protocol

# Example
```go
package main

import (

"fmt"
"github.com/fr05t1k/gelfkit"
"gopkg.in/Graylog2/go-gelf.v1/gelf"
)

func main() {
	gelfWriter, _ := gelf.NewWriter("localhost:12201")
	logger, _ := gelfkit.NewGelfLogger(gelfWriter)
	logger = log.With(logger, "caller", log.Caller(4))
	
	logger.Log("msg", "Hello world")
	
}
```

# Converting errors

You can covert `err` key to string by calling `EnableConvertErrors` method.
```go
    logger.EnableConvertErrors()
```

In this case if you pass an error in `err` field it will be converted by calling `Error()` method.
```go
    logger.Log("err", fmt.Errorf("test"))
```
It will  be converted to 
```
    {"err": "test"}
```