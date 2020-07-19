package podcast

// Item ..
type Item struct {
	Text        string     `xml:",chardata"`
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Encoded     string     `xml:"encoded"`
	Author      string     `xml:"author"`
	Summary     string     `xml:"summary"`
	Enclosure   *Enclosure `xml:"enclosure"`
	GUID        *GUID      `xml:"guid"`
	PubDate     string     `xml:"pubDate"`
	Duration    string     `xml:"duration"`
	Keywords    string     `xml:"keywords"`
	Season      string     `xml:"season"`
	Episode     string     `xml:"episode"`
	EpisodeType string     `xml:"episodeType"`
	Explicit    string     `xml:"explicit"`
}
