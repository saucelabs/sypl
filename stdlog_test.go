package sypl

import (
	"fmt"
	"log"
	"testing"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/output"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/safebuffer"
	"github.com/saucelabs/sypl/shared"
	"github.com/stretchr/testify/assert"
)

func getBufferLogger(name string, lvl level.Level) (*Sypl, *safebuffer.Buffer) {
	logger := NewDefault(name, lvl)
	buffer, outputBuffer := output.SafeBuffer(level.Trace, processor.PrefixBasedOnMask(shared.DefaultTimestampFormat))
	logger.AddOutputs(outputBuffer)

	return logger, buffer
}

func TestRedirectStdLogAt(t *testing.T) {
	initialFlags := log.Flags()
	initialPrefix := log.Prefix()

	levels := []level.Level{level.Trace, level.Debug, level.Info, level.Warn, level.Error}
	for _, lvl := range levels {
		runRedirectTest(t, lvl)
	}

	checkInitialFlags(t, initialFlags, initialPrefix)
}

func TestRedirectStdLogAtInvalid(t *testing.T) {
	logger := NewDefault("test-invalid", level.Debug)
	restore, err := RedirectStdLogAt(logger, level.None)
	defer func() {
		if restore != nil {
			restore()
		}
	}()

	assert.Error(t, err, "Expeted an error with an invalid log level")
}

func runRedirectTest(t *testing.T, lvl level.Level) {
	t.Helper()
	logger, logBuffer := getBufferLogger(fmt.Sprintf("test-%d", lvl), level.Trace)

	// test standard logger before and after redirect
	beforeMsg := "Before redirect"
	afterMsg := "After redirect"
	log.Println(beforeMsg)

	restore, err := RedirectStdLogAt(logger, lvl)
	assert.Nil(t, err, "Unexpected error redirecting std logs")
	defer restore()

	log.Println(afterMsg)
	checkLogMessage(t, beforeMsg, afterMsg, logBuffer.String())
}

func checkLogMessage(t *testing.T, beforeMsg string, afterMsg string, logBuffer string) {
	t.Helper()

	assert.NotContains(t, logBuffer, beforeMsg, "Log message written before redirect should not appear in Sypl logs")
	assert.Contains(t, logBuffer, afterMsg, "Log message written after redirect should appear in Sypl logs")
}

func checkInitialFlags(t *testing.T, initialFlags int, initialPrefix string) {
	t.Helper()
	assert.Equal(t, initialFlags, log.Flags(), "Expected to reset initial flags")
	assert.Equal(t, initialPrefix, log.Prefix(), "Expeted to reset initial prefix")
}
