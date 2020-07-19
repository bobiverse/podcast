package podcast

// Itunes constants
const (
	TypeEpisodic = "episodic"
	TypeSerial   = "serial"

	// ExplicitNo    = "no"
	// ExplicitYes   = "yes"
	ExplicitFalse = "false"
	ExplicitTrue  = "true"
)

// Type ..
func TypeValues() []string {
	return []string{
		TypeEpisodic,
		TypeSerial,
	}
}

// ExplicitValues ..
func ExplicitValues() []string {
	return []string{
		// ExplicitNo,
		// ExplicitYes,
		ExplicitFalse,
		ExplicitTrue,
	}
}
