package sypl

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/output"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/safebuffer"
	"github.com/saucelabs/sypl/shared"
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

	if (err == nil) {
		t.Error("Expected an error with an invalid log level")
	}
}

func runRedirectTest(t *testing.T, lvl level.Level) {
	logger, logBuffer := getBufferLogger(fmt.Sprintf("test-%d", lvl), level.Trace)

	// test standard logger before and after redirect
	beforeMsg := "Before redirect"
	afterMsg := "After redirect" 
	log.Println(beforeMsg)

	restore, err := RedirectStdLogAt(logger, lvl)
	if err != nil {
		t.Errorf("Unexpected error redirecting std logs: %x", err)
	}
	defer restore()

	log.Println(afterMsg)
	checkLogMessage(t, beforeMsg, afterMsg, logBuffer)

}

func checkLogMessage(t *testing.T, beforeMsg string, afterMsg string, logBuffer *safebuffer.Buffer) {

	bufferStr := logBuffer.String()

	if strings.Contains(bufferStr, beforeMsg) {
		t.Errorf("%s should not appear in proxy logger: %s", beforeMsg, bufferStr)
	}
	if !strings.Contains(bufferStr, afterMsg) {
		t.Errorf("%s should appear in proxy logger: %s", afterMsg, bufferStr)
	}
}

func checkInitialFlags(t *testing.T, initialFlags int, initialPrefix string) {
	if initialFlags != log.Flags() {
		t.Error("Expected to reset initial flags")
	}
	if initialPrefix != log.Prefix() {
		t.Error("Expected to reset initial prefix")
	}
}
