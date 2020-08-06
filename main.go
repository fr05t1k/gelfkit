package gelfkit

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"os"
	"path"
	"time"
)

const defaultLevel = gelf.LOG_INFO

type GelfWriter interface {
	WriteMessage(m *gelf.Message) (err error)
}

type GelfLogger struct {
	version       string
	gelfWriter    GelfWriter
	host          string
	facility      string
	convertErrors bool
}

func (l *GelfLogger) EnableConvertErrors() {
	l.convertErrors = true
}

func (l *GelfLogger) Log(keyvals ...interface{}) error {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, log.ErrMissingValue)
	}

	extra := map[string]interface{}{}

	for i := 0; i < len(keyvals)-1; i = i + 2 {
		key := keyvals[i]
		val := keyvals[i+1]
		extra[fmt.Sprintf("%v", key)] = val
	}

	// integration with "github.com/go-kit/kit/log/level"
	levelValue := getLevel(extra)
	delete(extra, "level")
	if l.convertErrors {
		extra = convertErrors(extra)
	}

	m := gelf.Message{
		Host:     l.host,
		Version:  l.version,
		TimeUnix: float64(time.Now().Unix()),
		Level:    levelValue,
		Facility: l.facility,
		Extra:    extra,
		RawExtra: nil,
	}

	return l.gelfWriter.WriteMessage(&m)
}

func convertErrors(extra map[string]interface{}) map[string]interface{} {
	if value, ok := extra["err"]; ok {
		if err, ok := value.(error); ok {
			extra["err"] = fmt.Sprintf("%s", err.Error())
		}
	}

	return extra
}

func getLevel(extra map[string]interface{}) int32 {
	value, ok := extra["level"]
	if !ok {
		return defaultLevel
	}

	levelValue, ok := value.(level.Value)
	switch levelValue {
	case level.ErrorValue():
		return gelf.LOG_ERR
	case level.DebugValue():
		return gelf.LOG_DEBUG
	case level.InfoValue():
		return gelf.LOG_INFO
	default:
		return defaultLevel
	}
}

func NewGelfLogger(writer GelfWriter) (*GelfLogger, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	l := &GelfLogger{
		gelfWriter:    writer,
		version:       "1.0",
		host:          host,
		facility:      path.Base(os.Args[0]),
		convertErrors: false,
	}

	return l, nil
}
