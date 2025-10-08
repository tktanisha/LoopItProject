package logger_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"loopit/pkg/logger"

	"github.com/golang/mock/gomock"
)

// helper: resets the singleton for each test
func resetLogger() {
	logger.ResetLoggerForTest()
}

func TestLogger_LogLevels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	logger.SetLogFileName(tmpFile.Name())

	resetLogger()
	l := logger.GetLogger()

	tests := []struct {
		name  string
		level func(msg string)
		msg   string
	}{
		{"Debug", l.Debug, "debug message"},
		{"Info", l.Info, "info message"},
		{"Warning", l.Warning, "warning message"},
		{"Error", l.Error, "error message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.level(tt.msg)
		})
	}

	time.Sleep(50 * time.Millisecond) // allow logs to flush

	l.Close()
	l.Close() // should not panic
}

func TestLogger_Fatal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmpFile, err := os.CreateTemp("", "fatal_log_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer os.Remove(tmpFile.Name())

	logger.SetLogFileName(tmpFile.Name())

	resetLogger()
	l := logger.GetLogger()

	exitCalled := false
	origExit := logger.OsExit
	logger.OsExit = func(code int) { exitCalled = true }
	defer func() { logger.OsExit = origExit }()
	l.Fatal("fatal message")
	if !exitCalled {
		t.Error("expected os.Exit to be called in Fatal()")
	}
}

func TestLogger_GetLogger_PanicOnOpenError(t *testing.T) {
	origOpenFile := logger.OpenFile
	logger.OpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("cannot open file")
	}
	defer func() { logger.OpenFile = origOpenFile }()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on open file failure")
		}
	}()
	resetLogger()
	logger.SetLogFileName("non_existent_path.log")
	logger.GetLogger()
}

func TestLogger_LogWhenClosing(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "closing_log_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	logger.SetLogFileName(tmpFile.Name())

	resetLogger()
	l := logger.GetLogger()
	l.Close()

	l.Debug("should not log") // should not panic
}
