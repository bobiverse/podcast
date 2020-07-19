package podcast

import "strings"

// Category  - https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12
type Category struct {
	// Text     string    `xml:",chardata"`
	AttrText string    `xml:"text,attr"`
	Category *Category `xml:"itunes:category"`
}

// IsEmpty ..
func (category *Category) IsEmpty() bool {
	return category == nil || category.AttrText == ""
}

// UnmarshalYAML ..
func (category *Category) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&category.AttrText)

	category.AttrText = strings.Trim(category.AttrText, " ,;/")
	arr := strings.Split(category.AttrText, ",")

	if len(arr) >= 1 {
		category.AttrText = strings.TrimSpace(arr[0])
	}

	if len(arr) >= 2 {
		category.Category = &Category{
			AttrText: strings.TrimSpace(arr[1]),
		}
	}

	return nil
}
