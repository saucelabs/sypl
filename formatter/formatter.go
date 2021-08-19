package formatter

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/processor"
	"github.com/saucelabs/sypl/shared"
)

// IFormatter specifies what a Formatter does.
type IFormatter = processor.IProcessor

//////
// Built-in processors.
//////

// JSON is a JSON formatter. It automatically adds:
// - Component name
// - Level
// - Timestamp (RFC3339).
func JSON() IFormatter {
	return processor.New("JSON", func(m message.IMessage) error {
		mM := map[string]interface{}{}

		mM["component"] = m.GetComponentName()
		mM["output"] = m.GetOutputName()
		mM["level"] = strings.ToLower(m.GetLevel().String())
		mM["timestamp"] = m.GetTimestamp().Format(time.RFC3339)
		mM["message"] = m.GetContent().GetProcessed()

		// Should only process fields if any.
		if len(m.GetFields()) != 0 {
			for k, v := range m.GetFields() {
				mM[k] = v
			}
		}

		m.GetContent().SetProcessed(shared.Prettify(mM))

		return nil
	})
}

// Text is a text formatter. It automatically adds:
// - Component name
// - Level
// - Timestamp (RFC3339).
func Text() IFormatter {
	return processor.New("Text", func(m message.IMessage) error {
		buf := new(strings.Builder)

		// Observe that the third line has no trailing tab,
		// so its final cell is not part of an aligned column.
		w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)

		fmt.Fprintf(w, "component=%s\t", m.GetComponentName())
		fmt.Fprintf(w, "output=%s\t", strings.ToLower(m.GetOutputName()))
		fmt.Fprintf(w, "level=%s\t", strings.ToLower(m.GetLevel().String()))
		fmt.Fprintf(w, "timestamp=%s\t", m.GetTimestamp().Format(time.RFC3339))
		fmt.Fprintf(w, "message=%s\t", m.GetContent().GetProcessed())

		// Should only process fields if any.
		if len(m.GetFields()) != 0 {
			for k, v := range m.GetFields() {
				fmt.Fprintf(w, "%s=%v\t", k, v)
			}
		}

		w.Flush()

		m.GetContent().SetProcessed(buf.String())

		return nil
	})
}
