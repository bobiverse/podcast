package podcast

// CDATA - <![CDATA[.. content
type CDATA struct {
	Text string `xml:",cdata"`
}

// UnmarshalYAML ..
func (cdata *CDATA) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&cdata.Text)
	return nil
}

// IsEmpty ..
func (cdata *CDATA) IsEmpty() bool {
	return cdata == nil || cdata.Text == ""
}

// String ..
func (cdata *CDATA) String() string {
	if cdata == nil {
		return ""
	}
	return cdata.Text
}
