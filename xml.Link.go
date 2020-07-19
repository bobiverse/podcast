package podcast

// Link ..
type Link struct {
	Text  string `xml:",chardata"`
	Href  string `xml:"href,attr"`
	Rel   string `xml:"rel,attr"`
	Type  string `xml:"type,attr"`
	Xmlns string `xml:"xmlns,attr"`
}
