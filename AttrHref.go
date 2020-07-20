package podcast

// AttrHref ..
type AttrHref struct {
	Href string `xml:"href,attr"`
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
