package helpers

import "testing"

func TestInitLogger(t *testing.T) {
	logger, err := InitLogger("DEBUG", true)
	if logger == nil {
		t.Error("logger is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	logger.Debug("Test")
}
