package podcast

import (
	"encoding/xml"
	"time"
)

// Date for feed better represent
type Date struct {
	time.Time
}

func (date Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := date.Time
	v := t.Format("2006-01-02")
	return e.EncodeElement(v, start)
}
