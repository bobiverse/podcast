package podcast

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// Duration for feed better represent
type Duration int

//
// func (dur Duration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// 	t := dur.Time
// 	// v := t.Format("2006-01-02")
// 	v := t.Format(time.RFC1123)
// 	return e.EncodeElement(v, start)
// }

// UnmarshalYAML ..
func (dur *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	unmarshal(&s)
	log.Printf(">> %s", s)

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
	log.Printf(">> %s", s)

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
