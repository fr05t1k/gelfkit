package gelfkit

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"os"
	"path"
	"time"
)

type GelfWriter interface {
	WriteMessage(m *gelf.Message) (err error)
}

type GelfLogger struct {
	version    string
	gelfWriter GelfWriter
	host       string
	facility   string
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

	m := gelf.Message{
		Host:     l.host,
		Version:  l.version,
		TimeUnix: float64(time.Now().Unix()),
		Level:    gelf.LOG_INFO,
		Facility: l.facility,
		Extra:    extra,
		RawExtra: nil,
	}

	return l.gelfWriter.WriteMessage(&m)
}

func NewGelfLogger(writer GelfWriter) (*GelfLogger, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	l := &GelfLogger{
		gelfWriter: writer,
		version:    "1.0",
		host:       host,
		facility:   path.Base(os.Args[0]),
	}

	return l, nil
}
