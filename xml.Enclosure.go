package podcast

// Enclosure ..
type Enclosure struct {
	// Text   string `xml:",chardata"`
	URL    string `xml:"url,omitempty,attr"`
	Length int64  `xml:"length,omitempty,attr"`
	Type   string `xml:"type,omitempty,attr"`
}

// IsEmpty ..
func (enc *Enclosure) IsEmpty() bool {
	return enc.URL == "" || enc.Length == 0 || enc.Type == ""
}
