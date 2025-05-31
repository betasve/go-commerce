package jsonlog

import (
	"bytes"
	"encoding/json"
	"errors"
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
	var payload logBody

	tests := []struct {
		name            string
		logger          *Logger
		expectedLevel   string
		expectedMessage string
		expectedPropKey string
		expectedPropVal string
	}{
		{"Prints INFO to the log", New(&buf, LevelInfo), "INFO", "info msg", "err", "val"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.logger.PrintInfo(tc.expectedMessage, map[string]string{tc.expectedPropKey: tc.expectedPropVal})

			err := json.NewDecoder(&buf).Decode(&payload)
			if err != nil {
				t.Error(err.Error())
			}

			if payload.Level != tc.expectedLevel {
				t.Errorf("Expected %v and got %v", tc.expectedLevel, payload.Level)
			}

			if payload.Message != tc.expectedMessage {
				t.Errorf("Expected %v and got %v", tc.expectedMessage, payload.Message)
			}

			if payload.Properties[tc.expectedPropKey] != tc.expectedPropVal {
				t.Errorf("Expected %v and got %v", tc.expectedPropVal, payload.Properties[tc.expectedPropKey])
			}
		})
	}
}

func TestPrintFatal(t *testing.T) {
	var buf bytes.Buffer
	var payload logBody
	var exitCalled bool
	log := &Logger{
		out:      &buf,
		minLevel: LevelFatal,
		exitFn:   func(code int) { exitCalled = true },
	}

	log.PrintFatal(errors.New("some fatal error"), map[string]string{"key": "fatal error"})

	err := json.NewDecoder(&buf).Decode(&payload)
	if err != nil {
		t.Error(err.Error())
	}

	if payload.Level != "FATAL" {
		t.Errorf("Expected %v and got %v", "FATAL", payload.Level)
	}

	if payload.Message != "some fatal error" {
		t.Errorf("Expected %v and got %v", "some fatal error", payload.Message)
	}

	if payload.Properties["key"] != "fatal error" {
		t.Errorf("Expected %v and got %v", "fatal error", payload.Properties["key"])
	}

	if !exitCalled {
		t.Errorf("Expected to call the exit function but it didn't.")
	}
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	var payload logBody

	tests := []struct {
		name            string
		logger          *Logger
		expectedLevel   string
		expectedMessage error
		expectedPropKey string
		expectedPropVal string
	}{
		{"Prints ERROR to the log", New(&buf, LevelInfo), "ERROR", errors.New("test"), "err", "val"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.logger.PrintError(tc.expectedMessage, map[string]string{tc.expectedPropKey: tc.expectedPropVal})

			err := json.NewDecoder(&buf).Decode(&payload)
			if err != nil {
				t.Error(err.Error())
			}

			if payload.Level != tc.expectedLevel {
				t.Errorf("Expected %v and got %v", tc.expectedLevel, payload.Level)
			}

			if payload.Message != tc.expectedMessage.Error() {
				t.Errorf("Expected %v and got %v", tc.expectedMessage, payload.Message)
			}

			if payload.Properties[tc.expectedPropKey] != tc.expectedPropVal {
				t.Errorf("Expected %v and got %v", tc.expectedPropVal, payload.Properties[tc.expectedPropKey])
			}
		})
	}
}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	var payload logBody

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
		{"Info level gets logged", New(&buf, LevelInfo), LevelInfo, "info msg", map[string]string{"key": "err"}, "info msg", "INFO", "err", true},
		{"Error level gets logged", New(&buf, LevelError), LevelError, "error msg", map[string]string{"key": "err2"}, "error msg", "ERROR", "err2", true},
		{"Fatal level gets logged", New(&buf, LevelFatal), LevelFatal, "fatal msg", map[string]string{"key": "err3"}, "fatal msg", "FATAL", "err3", true},
		{"Error level with Info does not log", New(&buf, LevelError), LevelInfo, "fatal msg", map[string]string{"key": "err3"}, "fatal msg", "FATAL", "err3", false},
		{"Fatal level with Info does not log", New(&buf, LevelFatal), LevelInfo, "fatal msg", map[string]string{"key": "err3"}, "fatal msg", "FATAL", "err3", false},
		{"Fatal level with Error does not log", New(&buf, LevelFatal), LevelError, "fatal msg", map[string]string{"key": "err3"}, "fatal msg", "FATAL", "err3", false},
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
