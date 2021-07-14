package flag

// Flag defines behaviours.
type Flag int

const (
	None Flag = iota

	// Force message will be processed, and printed independent of `Level`
	// restrictions.
	Force

	// Mute message will be processed, but not printed.
	Mute

	// Skip message will not be processed.
	Skip

	// SkipAndForce message will not be processed, but will be printed
	// independent of `Level` restrictions.
	SkipAndForce

	// SkipAndMute message will not be processed, neither printed.
	SkipAndMute
)

var names = [...]string{"None", "Force", "Mute", "Skip"}

// String translates enum flags to string.
func (f Flag) String() string {
	if f < None || f > Skip {
		return "Unknown"
	}

	return names[f]
}
