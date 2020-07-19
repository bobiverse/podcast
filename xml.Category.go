package podcast

// Category ..
type Category struct {
	Text     string    `xml:",chardata"`
	AttrText string    `xml:"text,attr"`
	Category *Category `xml:"category"`
}
