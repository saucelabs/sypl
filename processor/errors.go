// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package processor

import (
	"fmt"
	"strings"

	"github.com/saucelabs/sypl/message"
)

// ProcessingError is thrown when a processor fails.
type ProcessingError struct {
	// Cause is the underlying case of the failure.
	Cause error

	// Message when the failured happened.
	Message message.IMessage

	// OutputName is the name of the output in-use.
	OutputName string

	// ProcessorName is the name of the processor in-use.
	ProcessorName string
}

// Error interface implementation.
func (p *ProcessingError) Error() string {
	errMsg := fmt.Sprintf(`Output: "%s" Processor: "%s"`,
		p.OutputName,
		p.ProcessorName,
	)

	if p.Cause != nil {
		errMsg = fmt.Sprintf(`%s Error: "%s"`, errMsg, p.Cause)
	}

	if p.Message != nil {
		errMsg = fmt.Sprintf(`%s Original Message: "%s"`,
			errMsg,
			strings.TrimSuffix(p.Message.GetContent().GetOriginal(), "\n"),
		)
	}

	return errMsg
}

// NewProcessingError returns a new `ProcessingError`.
func NewProcessingError(m message.IMessage, e error) error {
	return &ProcessingError{
		Cause:         e,
		Message:       m,
		OutputName:    m.GetOutputName(),
		ProcessorName: m.GetProcessorName(),
	}
}
