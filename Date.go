package podcast

import (
	"encoding/xml"
	"time"
)

// Date for feed better represent
type Date struct {
	time.Time
}

// MarshalXML ..
func (date Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := date.Time
	// v := t.Format("2006-01-02")
	v := t.Format(time.RFC1123) // Wed, 22 Jul 2020 13:01:20 EEST
	// v := t.Format(time.RFC822) // 01 Jan 20 11:59 UTC
	return e.EncodeElement(v, start)
}

// IsZero ..
func (date *Date) IsZero() bool {
	return date == nil || date.Time.IsZero()
}

// UnmarshalYAML ..
func (date *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&date.Time)

	return nil
}
