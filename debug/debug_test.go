package debug

import (
	"os"
	"testing"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/shared"
)

func TestNew(t *testing.T) {
	type args struct {
		componentName string
		outputName    string
	}
	tests := []struct {
		name        string
		args        args
		content     string
		wantLevel   level.Level
		wantMatcher Matcher
		wantOK      bool
	}{
		{
			name: "Should work - COL",
			args: args{
				componentName: "componentX",
				outputName:    "outputY",
			},
			content:     "info,componentX:outputY:debug,outputZ:trace",
			wantLevel:   level.Debug,
			wantMatcher: COL,
			wantOK:      true,
		},
		{
			name: "Should work - OL",
			args: args{
				componentName: "componentX",
				outputName:    "outputZ",
			},
			content:     "info,componentX:outputY:debug,outputZ:trace",
			wantLevel:   level.Trace,
			wantMatcher: OL,
			wantOK:      true,
		},
		{
			name: "Should work - L",
			args: args{
				componentName: "componentX",
				outputName:    "outputW",
			},
			content:     "info,componentX:outputY:debug,outputZ:trace",
			wantLevel:   level.Info,
			wantMatcher: L,
			wantOK:      true,
		},
		{
			name: "Should work - level none",
			args: args{
				componentName: "componentX",
				outputName:    "outputY",
			},
			content:     "info,componentX:outputY:none,outputZ:trace",
			wantLevel:   level.None,
			wantMatcher: COL,
			wantOK:      true,
		},
		{
			name: "Should fail - no match",
			args: args{
				componentName: "componentX",
				outputName:    "outputW",
			},
			content:     "componentX:outputY:debug,outputZ:trace,info",
			wantLevel:   level.None,
			wantMatcher: None,
			wantOK:      false,
		},
		{
			name: "Should fail - invalid level",
			args: args{
				componentName: "componentX",
				outputName:    "outputY",
			},
			content:     "componentX:outputY:asd,outputZ:trace,info",
			wantLevel:   level.None,
			wantMatcher: None,
			wantOK:      false,
		},
		{
			name: "Should do nothing",
			args: args{
				componentName: "componentX",
				outputName:    "outputY",
			},
			wantLevel:   level.None,
			wantMatcher: None,
			wantOK:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(shared.DebugEnvVar, tt.content)
			defer os.Unsetenv(shared.DebugEnvVar)

			d := New(tt.args.componentName, tt.args.outputName)

			lvl, m, ok := d.Level()

			if ok != tt.wantOK {
				t.Fatalf("OK; expected %+v got %+v", tt.wantOK, ok)
			}

			if m != tt.wantMatcher {
				t.Fatalf("OK; expected %+v got %+v", tt.wantMatcher, m)
			}

			if lvl != tt.wantLevel {
				t.Fatalf("OK; expected %+v got %+v", tt.wantLevel, lvl)
			}
		})
	}
}
