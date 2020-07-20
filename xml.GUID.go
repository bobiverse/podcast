package podcast

// GUID ..
type GUID struct {
	Text        string `xml:",chardata"`
	IsPermaLink bool   `xml:"isPermaLink,attr"`
}

// IsEmpty ..
func (guid *GUID) IsEmpty() bool {
	return guid == nil || guid.Text == ""
}

// UnmarshalYAML ..
func (guid *GUID) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&guid.Text)

	guid.IsPermaLink = isValidURL(guid.Text)

	return nil
}

// NewGUID ..
func NewGUID(s string) *GUID {
	guid := &GUID{
		Text:        s,
		IsPermaLink: isValidURL(s),
	}
	return guid
}
