package logs

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	var logger = GetLogger("WARN")
	logger.Infof("testing INFO level")
}
