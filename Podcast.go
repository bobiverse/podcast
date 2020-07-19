package podcast

import (
	"log"

	"github.com/fatih/color"
)

// Podcast ..
type Podcast struct {
	ContentPath string
	OutputPath  string
	Feed        *XMLRoot
}

// New ..
func New() *Podcast {
	log.Println("[podcast] New ")

	podcast := &Podcast{
		Feed: &XMLRoot{
			Itunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
			Content: "http://purl.org/rss/1.0/modules/content/",
			Atom:    "http://www.w3.org/2005/Atom",
			Version: "2.0",
		},
	}
	return podcast
}

// SaveToFile ..
func (podcast *Podcast) SaveToFile() error {
	log.Println("[podcast] SaveToFile ")

	buf, err := podcast.Feed.ToXML("")
	if err != nil {
		return err
	}

	color.Yellow("%s", buf)
	return nil
}
