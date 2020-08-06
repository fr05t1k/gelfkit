package gelfkit

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
	"reflect"
	"testing"
)

type HistoryGelfWriter struct {
	history []*gelf.Message
}

func (h *HistoryGelfWriter) WriteMessage(m *gelf.Message) (err error) {
	h.history = append(h.history, m)
	return
}

func (h *HistoryGelfWriter) GetLastMessage() *gelf.Message {
	if len(h.history) > 0 {
		return h.history[len(h.history)-1]
	}

	return nil
}

func NewHistoryWriter() *HistoryGelfWriter {
	return &HistoryGelfWriter{history: []*gelf.Message{}}
}

func TestGelfLogger_Log(t *testing.T) {
	tests := []struct {
		name                  string
		keyvals               []interface{}
		enableErrorConverting bool
		want                  gelf.Message
		wantErr               bool
	}{
		{
			name: "simple",
			keyvals: []interface{}{
				"test", "test",
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "num",
			keyvals: []interface{}{
				123, 123,
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"123": 123,
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate",
			keyvals: []interface{}{
				0, 123,
				0, 456,
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"0": 456,
				},
			},
			wantErr: false,
		},
		{
			name: "complex key",
			keyvals: []interface{}{
				map[string]int{"test": 1}, 123,
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"map[test:1]": 123,
				},
			},
			wantErr: false,
		},
		{
			name: "odd numbers of keys",
			keyvals: []interface{}{
				"name", "value",
				"name2",
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"name":  "value",
					"name2": log.ErrMissingValue,
				},
			},
			wantErr: false,
		},
		{
			name: "integration with level package",
			keyvals: []interface{}{
				level.Key(), level.DebugValue(),
				"key", "value",
			},
			want: gelf.Message{
				Level: gelf.LOG_DEBUG,
				Extra: map[string]interface{}{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name:                  "error converting",
			enableErrorConverting: true,
			keyvals: []interface{}{
				"err", fmt.Errorf("test"),
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"err": "test",
				},
			},
			wantErr: false,
		},
		{
			name:                  "disable error converting",
			enableErrorConverting: false,
			keyvals: []interface{}{
				"err", fmt.Errorf("test"),
			},
			want: gelf.Message{
				Extra: map[string]interface{}{
					"err": fmt.Errorf("test"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewHistoryWriter()
			l, err := NewGelfLogger(w)
			if err != nil {
				t.Errorf("NewGelfLogger() error = %v", err)
			}
			if tt.enableErrorConverting {
				l.EnableConvertErrors()
			}

			if err := l.Log(tt.keyvals...); (err != nil) != tt.wantErr {
				t.Errorf("Log() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := w.GetLastMessage()
			if got == nil {
				t.Errorf("Log() = nil, want %v", tt.want)
				return
			}

			if !reflect.DeepEqual(tt.want.Extra, got.Extra) {
				t.Errorf("Log() = %v, want %v", got.Extra, tt.want.Extra)
			}
		})
	}
}

func Test_getLevel(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]interface{}
		want  int32
	}{
		{
			name:  "info",
			extra: map[string]interface{}{"level": level.InfoValue()},
			want:  gelf.LOG_INFO,
		},
		{
			name:  "error",
			extra: map[string]interface{}{"level": level.ErrorValue()},
			want:  gelf.LOG_ERR,
		},
		{
			name:  "debug",
			extra: map[string]interface{}{"level": level.DebugValue()},
			want:  gelf.LOG_DEBUG,
		},
		{
			name:  "no key",
			extra: map[string]interface{}{"level2": level.DebugValue()},
			want:  defaultLevel,
		},
		{
			name:  "wrong type",
			extra: map[string]interface{}{"level": "debug"},
			want:  defaultLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLevel(tt.extra); got != tt.want {
				t.Errorf("getLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
