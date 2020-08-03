package logger_test

import (
	"testing"

	"github.com/belldata/bdx/logger"
)

func logging(log logger.ILogger) {
	log.Debug("debug")
	log.Info("Info")
	log.Warning("Warning")
	log.Error("Error")
	log.Fatal("Fatal")

	log.Debugf("format: %s", "debug")
	log.Infof("format: %s", "Info")
	log.Warningf("format: %s", "Warning")
	log.Errorf("format: %s", "Error")
	log.Fatalf("format: %s", "Fatal")
}

func TestBxLogger(t *testing.T) {
	log := logger.New("test", logger.Debug)
	// logger.SetLogColor(true)
	logging(log)
	println()
	log.SetLevel(logger.Info)
	logging(log)
	println()
	log.SetLevel(logger.Warning)
	logging(log)
	println()
	log.SetLevel(logger.Error)
	logging(log)
	println()
	log.SetLevel(logger.Fatal)
	logging(log)
}
