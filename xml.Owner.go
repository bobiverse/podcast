package podcast

import "strings"

// Owner ..
type Owner struct {
	// Text  string `xml:",chardata"`
	Name  string `xml:"itunes:name,omitempty"`
	Email string `xml:"itunes:email,omitempty"`
}

// IsEmpty ..
func (owner *Owner) IsEmpty() bool {
	return owner == nil || owner.Name == "" || owner.Email == ""
}

// UnmarshalYAML ..
func (owner *Owner) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&owner.Name)

	s := strings.Trim(owner.Name, " ,;/")
	arr := strings.Split(s, ",")

	if len(arr) >= 1 {
		owner.Name = strings.TrimSpace(arr[0])
	}

	if len(arr) >= 2 {
		owner.Email = strings.TrimSpace(arr[1])
	}

	return nil
}
