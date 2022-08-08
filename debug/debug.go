package debug

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/saucelabs/sypl/level"
	"github.com/saucelabs/sypl/shared"
)

type Matcher string

const (
	// L matches against any valid level, specified at the beginning of the
	// debug env var, examples:
	// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `info`
	// - SYPL_DEBUG="componentX:outputY:debug,outputZ:trace,info" -> ``.
	//
	// Note: For this matcher, the order matter!
	L Matcher = "Level"

	// OL matches against a specific output, and any valid level specified in
	// the debug env var, example:
	// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `trace`.
	OL Matcher = "OutputLevel"

	// COL matches against a specific component and output, and any valid level
	// specified in the debug env var, example:
	// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `debug`.
	COL Matcher = "ComponentOutputLevel"

	// None means no Matcher matched against the debug env var.
	None Matcher = "None"
)

// Matchers' regexes.
const (
	lReMask   = `(?i)^(?:%s)`
	oLReMask  = `(?i)(?:(?:%s):(?:%s))`
	cOLReMask = `(?i)(?:(?:%s):(?:%s):(?:%s))`
)

// Debug definition.
type Debug struct {
	// ComponentName is the component name.
	ComponentName string

	// OutputName is the output name.
	OutputName string

	// Content of the debug env var.
	Content string

	// Levels matcher regex matches against any valid level, specified at the
	// beginning of the debug env var, examples:
	// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `info`
	// - SYPL_DEBUG="componentX:outputY:debug,outputZ:trace,info" -> ``.
	//
	// Note: For this matcher, the order matter!
	Levels *regexp.Regexp

	// Output, and levels matcher regex matches against a specific output, and
	// any valid level specified in the debug env var, example:
	// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `trace`
	OutputLevels *regexp.Regexp

	// COL matches against a specific component and output, and any valid level
	// specified in the debug env var, example:
	// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `debug`.
	ComponentOutputLevels *regexp.Regexp
}

// MatchL uses the `Levels` matcher against any valid level, specified at the
// beginning of the debug env var, examples:
// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `info`
// - SYPL_DEBUG="componentX:outputY:debug,outputZ:trace,info" -> “.
//
// Notes:
// - For this matcher, the order matter!
// - Prefer to use the `Level` method.
func (d *Debug) MatchL() string {
	return d.Levels.FindString(d.Content)
}

// MatchOL uses the `OutputLevels` matcher against a specific output, and
// any valid level specified in the debug env var, example:
// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `trace`
//
// Note: Prefer to use the `Level` method.
func (d *Debug) MatchOL() string {
	return d.OutputLevels.FindString(d.Content)
}

// MatchCOL uses the `ComponentOutputLevels` matcher against a specific
// component and output, and any valid level specified in the debug env var,
// example:
// - SYPL_DEBUG="info,componentX:outputY:debug,outputZ:trace" -> `debug`.
//
// Note: Prefer to use the `Level` method.
func (d *Debug) MatchCOL() string {
	return d.ComponentOutputLevels.FindString(d.Content)
}

// Level checks the content of the debug env var against all matchers returning:
// - The level extracted from the last Matcher
// - The last `Matcher` that matched
// - If any matcher succeeded on matching
//
// Matchers:
// - {componentName:outputName:level} -> forwarder:console:trace
// - {outputName:level} -> console:trace
// - {level}, e.g.: trace
//
// Note: Don't use the returned level to check if `Level` succeeded because
// `level.None` is a valid, and usable level.
func (d *Debug) Level() (level.Level, Matcher, bool) {
	// Shouldn't' do anything if the debug env var isn't set.
	if d.Content == "" {
		return level.None, None, false
	}

	var (
		finalMatcher       Matcher
		finalLevelAsString string
	)

	// Matches' order matters.
	if lReMatch := d.MatchL(); lReMatch != "" {
		finalLevelAsString = strings.Split(lReMatch, ":")[0]
		finalMatcher = L
	}

	if oLReMatch := d.MatchOL(); oLReMatch != "" {
		finalLevelAsString = strings.Split(oLReMatch, ":")[1]
		finalMatcher = OL
	}

	if nOLReMatch := d.MatchCOL(); nOLReMatch != "" {
		finalLevelAsString = strings.Split(nOLReMatch, ":")[2]
		finalMatcher = COL
	}

	finalLevel, err := level.FromString(finalLevelAsString)
	if err != nil {
		return level.None, None, false
	}

	return finalLevel, finalMatcher, true
}

//////
// Factory.
//////

// New is the Debug factory.
func New(componentName, outputName string) *Debug {
	levels := strings.Join(level.LevelsNames(), "|")

	return &Debug{
		ComponentName: componentName,
		OutputName:    outputName,

		Content: os.Getenv(shared.DebugEnvVar),

		Levels:                regexp.MustCompile(fmt.Sprintf(lReMask, levels)),
		OutputLevels:          regexp.MustCompile(fmt.Sprintf(oLReMask, outputName, levels)),
		ComponentOutputLevels: regexp.MustCompile(fmt.Sprintf(cOLReMask, componentName, outputName, levels)),
	}
}
