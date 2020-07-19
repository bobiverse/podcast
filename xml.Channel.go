package podcast

// Channel ..
type Channel struct {
	Text          string    `xml:",chardata"`
	Link          *Link     `xml:"link"`
	Title         string    `xml:"title"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Language      string    `xml:"language"`
	Copyright     string    `xml:"copyright"`
	Subtitle      string    `xml:"subtitle"`
	Author        string    `xml:"author"`
	Summary       string    `xml:"summary"`
	Type          string    `xml:"type"`
	Explicit      string    `xml:"explicit"`
	Description   string    `xml:"description"`
	Keywords      string    `xml:"keywords"`
	Owner         *Owner    `xml:"owner"`
	Image         *Image    `xml:"image"`
	Category      *Category `xml:"category"`
	Items         ItemList  `xml:"item"`
}
