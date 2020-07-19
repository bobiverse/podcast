package podcast

// Image ..
type Image struct {
	Text  string `xml:",chardata"`
	Href  string `xml:"href,attr"`
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}
