package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/shared"
)

// IFormatter specifies what a Formatter does.
type IFormatter = processor.IProcessor

// JSON is a JSON formatter. It automatically adds:
// - Component name
// - Level
// - Timestamp (RFC3339).
func JSON() IFormatter {
	return processor.NewProcessor("JSON", func(m message.IMessage) error {
		mM := map[string]interface{}{}

		mM["component"] = m.GetComponentName()
		mM["content"] = m.GetContent().GetProcessed()
		mM["level"] = strings.ToLower(m.GetLevel().String())

		// Should only process fields if any.
		if len(m.GetFields()) != 0 {
			for k, v := range m.GetFields() {
				mM[k] = v
			}
		}

		mM["timestamp"] = m.GetTimestamp().Format(time.RFC3339)

		m.GetContent().SetProcessed(shared.Prettify(mM))

		return nil
	})
}

// Text is a text formatter. It automatically adds:
// - Component name
// - Level
// - Timestamp (RFC3339).
func Text() IFormatter {
	return processor.NewProcessor("Text", func(m message.IMessage) error {
		finalMessage := m.GetContent().GetProcessed()

		buf := new(strings.Builder)

		fmt.Fprintf(buf, "component=%v ", m.GetComponentName())
		fmt.Fprintf(buf, "level=%v ", strings.ToLower(m.GetLevel().String()))

		// Should only process fields if any.
		if len(m.GetFields()) != 0 {
			for k, v := range m.GetFields() {
				fmt.Fprintf(buf, "%s=%v ", k, v)
			}
		}

		fmt.Fprintf(buf, "timestamp=%v ", m.GetTimestamp().Format(time.RFC3339))

		processedField := strings.TrimSuffix(buf.String(), " ")

		finalMessage = fmt.Sprintf("%s %s", finalMessage, processedField)

		m.GetContent().SetProcessed(finalMessage)

		return nil
	})
}
