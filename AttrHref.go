package podcast

// AttrHref ..
type AttrHref struct {
	Href string `xml:"href,omitempty,attr"`
	Rel  string `xml:"rel,omitempty,attr"`
	Type string `xml:"type,omitempty,attr"`
}

// IsEmpty ..
func (href *AttrHref) IsEmpty() bool {
	return href == nil || href.Href == ""
}

// UnmarshalYAML ..
func (href *AttrHref) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&href.Href)
	return nil
}
