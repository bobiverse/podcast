package podcast

// Itunes constants
const (
	PodcastTypeEpisodic = "episodic"
	PodcastTypeSerial   = "serial"

	EpisodeTypeFull    = "full"
	EpisodeTypeTrailer = "trailer"
	EpisodeTypeBonus   = "bonus"

	ExplicitNo    = "no"
	ExplicitYes   = "yes"
	ExplicitFalse = "false"
	ExplicitTrue  = "true"
)

// PodcastTypeValues ..
func PodcastTypeValues() []string {
	return []string{
		PodcastTypeEpisodic,
		PodcastTypeSerial,
	}
}

// EpisodeTypesValues ..
func EpisodeTypesValues() []string {
	return []string{
		EpisodeTypeFull,
		EpisodeTypeTrailer,
		EpisodeTypeBonus,
	}
}

// ExplicitValues ..
func ExplicitValues() []string {
	return []string{
		ExplicitNo,
		ExplicitYes,
		ExplicitFalse,
		ExplicitTrue,
	}
}
