package podcast

import (
	"strconv"
	"strings"
	"time"
)

// Duration for feed better represent
// From apple: Different duration formats are accepted however it is recommended to convert the length of the episode into seconds.
type Duration int

// UnmarshalYAML ..
func (dur *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	unmarshal(&s)

	if err := dur.Set(s); err != nil {
		return err
	}

	return nil
}

// Set duration from string
func (dur *Duration) Set(s string) error {
	// Try to convert to integer if no errors then ok
	totalSeconds, err := strconv.Atoi(s)
	if err == nil {
		*dur = Duration(totalSeconds)
		return nil
	}

	// If not try to convert from common duration formats

	// MM:SS ==> HH:MM:SS
	// Add missing hours
	if strings.Count(s, ":") == 1 && len(s) == 5 {
		s = "00:" + s
	}

	// correct to go format for duration parsing
	s = strings.Replace(s, ":", "h", 1) // 00:52:11 ==> 00h52:11
	s = strings.Replace(s, ":", "m", 1) // 00h52:11 ==>  00h52m11
	s += "s"                            // 00h52m11s
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*dur = Duration(d.Seconds())
	return nil
}
