package formatter

import (
	"strings"
	"testing"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/message"
	"github.com/saucelabs/sypl/options"
	"github.com/saucelabs/sypl/shared"
)

func TestText(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := message.NewMessage(level.Info, shared.DefaultContentOutput)
			m.SetComponentName(shared.DefaultComponentNameOutput)
			m.SetFields(options.Fields{
				"key1": "value1",
			})

			if err := Text().Run(m); err != nil {
				t.Errorf("Text() = %v, error %v", m, err)
			}

			if !strings.Contains(m.String(), "component=") {
				t.Errorf("Text() = missing %s", "component=")
			}
			if !strings.Contains(m.String(), shared.DefaultContentOutput) {
				t.Errorf("Text() = missing %s", shared.DefaultContentOutput)
			}
			if !strings.Contains(m.String(), `key1=value1`) {
				t.Errorf("Text() = missing %s", `key1=value1`)
			}
			if !strings.Contains(m.String(), "level=") {
				t.Errorf("Text() = missing %s", "level=")
			}
			if !strings.Contains(m.String(), "timestamp=") {
				t.Errorf("Text() = missing %s", "timestamp=")
			}
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := message.NewMessage(level.Info, shared.DefaultContentOutput)
			m.SetComponentName(shared.DefaultComponentNameOutput)
			m.SetFields(options.Fields{
				"key1": "value1",
			})

			if err := JSON().Run(m); err != nil {
				t.Errorf("Text() = %v, error %v", m, err)
			}

			if !strings.Contains(m.String(), `"component"`) {
				t.Errorf("Text() = missing %s", `"component"`)
			}
			if !strings.Contains(m.String(), shared.DefaultContentOutput) {
				t.Errorf("Text() = missing %s", shared.DefaultContentOutput)
			}
			if !strings.Contains(m.String(), `"key1": "value1"`) {
				t.Errorf("Text() = missing %s", `"key1": "value1"`)
			}
			if !strings.Contains(m.String(), `"level"`) {
				t.Errorf("Text() = missing %s", `"level"`)
			}
			if !strings.Contains(m.String(), `"timestamp"`) {
				t.Errorf("Text() = missing %s", `"timestamp"`)
			}
		})
	}
}
