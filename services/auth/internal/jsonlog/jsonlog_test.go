package jsonlog

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name          string
		level         Level
		expectedLevel string
	}{
		{"Info level", LevelInfo, "INFO"},
		{"Error level", LevelError, "ERROR"},
		{"Fatal level", LevelFatal, "FATAL"},
		{"Off", LevelOff, "OFF"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.level.String()
			if result != tc.expectedLevel {
				t.Errorf("Expected '%v', got '%v'", tc.expectedLevel, result)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	expectedLevel := Level(1)

	logger := New(&buf, expectedLevel)

	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}

	if logger.out != &buf {
		t.Errorf("expected out to be %v, got %v", &buf, logger.out)
	}

	if logger.minLevel != expectedLevel {
		t.Errorf("expected minLevel to be %v, got %v", expectedLevel, logger.minLevel)
	}
}

func TestPrintInfo(t *testing.T) {
	var buf bytes.Buffer

	var payload struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}

	tests := []struct {
		name            string
		logger          *Logger
		level           Level
		message         string
		props           map[string]string
		expectedMessage string
		expectedLevel   string
		expectedVal     string
		logged          bool
	}{
		{"Info level", New(&buf, LevelInfo), LevelInfo, "info msg", map[string]string{"key": "err"}, "info msg", "INFO", "err", true},
		{"Error level", New(&buf, LevelError), LevelError, "error msg", map[string]string{"key": "err2"}, "error msg", "ERROR", "err2", true},
		{"Fatal level", New(&buf, LevelFatal), LevelFatal, "fatal msg", map[string]string{"key": "err3"}, "fatal msg", "FATAL", "err3", true},
		{"Error level with info", New(&buf, LevelError), LevelInfo, "fatal msg", map[string]string{"key": "err3"}, "fatal msg", "FATAL", "err3", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.logger.print(tc.level, tc.message, tc.props)
			if err != nil {
				t.Error(err.Error())
			}

			if !tc.logged {
				_, err := buf.ReadBytes('0')
				if err != io.EOF {
					t.Errorf("Expected EOF and got %v", err)
				}

				return
			}

			err = json.NewDecoder(&buf).Decode(&payload)
			if err != nil {
				t.Error(err.Error())
			}

			if payload.Level != tc.expectedLevel {
				t.Errorf("Expected %v and got %v", tc.expectedLevel, payload.Level)
			}

			if payload.Message != tc.expectedMessage {
				t.Errorf("Expected %v and got %v", tc.expectedMessage, payload.Message)
			}

			if payload.Properties["key"] != tc.expectedVal {
				t.Errorf("Expected %v and got %v", tc.expectedVal, payload.Properties["key"])
			}
		})
	}
}
