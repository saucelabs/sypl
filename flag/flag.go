package flag

// Flag defines behaviours.
type Flag int

const (
	None Flag = iota

	// Force if set message will be processed, and printed independent of
	// `Level` restrictions.
	Force

	// Mute if set, message will not be printed.
	Mute

	// Skip if set, processors will ignore the message, and do nothing.
	Skip
)

var names = [...]string{"None", "Force", "Mute", "Skip"}

// String translates enum flags to string.
func (f Flag) String() string {
	if f < None || f > Skip {
		return "Unknown"
	}

	return names[f]
}
