package podcast

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Item ..
type Item struct {
	Key      string `xml:"-" yaml:"-"`
	File     string `xml:"-" yaml:"File"`
	FileSize int64 `xml:"-" yaml:"FileSize"`
	FileType string `xml:"-" yaml:"FileType"`

	// Text        string     `xml:",chardata"`
	Title       string     `xml:"title,omitempty" yaml:"Title"`
	Description string     `xml:"description,omitempty" yaml:"Description"`
	Encoded     *CDATA     `xml:"content:encoded,omitempty" yaml:"Encoded"`
	Author      string     `xml:"author,omitempty" yaml:"Author"`
	Summary     *CDATA     `xml:"summary,omitempty" yaml:"Summary"`
	Enclosure   *Enclosure `xml:"enclosure,omitempty" yaml:"Enclosure"`
	GUID        *GUID      `xml:"guid,omitempty" yaml:"GUID"`
	PubDate     *Date      `xml:"pubDate,omitempty" yaml:"PubDate"`
	Duration    string     `xml:"itunes:duration,omitempty" yaml:"Duration"`
	Keywords    string     `xml:"itunes:keywords,omitempty" yaml:"Keywords"`
	Season      int        `xml:"itunes:season,omitempty" yaml:"Season"`
	Episode     int        `xml:"itunes:episode,omitempty" yaml:"Episode"`
	EpisodeType string     `xml:"itunes:episodeType,omitempty" yaml:"EpisodeType"`
	Explicit    string     `xml:"itunes:explicit,omitempty" yaml:"Explicit"`

	ItunesTitle   string `xml:"itunes:title,omitempty" yaml:"ItunesTitle"`
	ItunesSummary *CDATA `xml:"itunes:summary,omitempty" yaml:"ItunesSummary"`
}

// Weight of the item for sorting
// Seasons and episode taken
func (item *Item) Weight() int {
	weight := item.Season*1000 + item.Episode
	log.Printf("WEIGHT: %d", weight)
	return weight
}

// Extraxt info from key
// S01E02 ..
func (item *Item) ExtractKeyInfo() {
	if len(item.Key) != 6 {
		return
	}

	// Season
	season, err := strconv.Atoi(item.Key[1:3])
	if err != nil {
		log.Printf("Warning: Season can't be extracted from Item [%s]", item.Key)

	} else if item.Season == 0 {
		// assign if not assigned in yaml file
		item.Season = season
	}

	// Episode
	episode, err := strconv.Atoi(item.Key[5:6])
	if err != nil {
		log.Printf("Warning: Episode can't be extracted from Item [%s]", item.Key)

	} else if item.Episode == 0 {
		// assign if not assigned in yaml file
		item.Episode = episode
	}
}

// Fix ..
func (item *Item) Fix() {
	log.Printf("Item[%s] Fix()...", item.Key)

	if item.ItunesTitle == "" {
		item.ItunesTitle = item.Title
	}

	if item.ItunesSummary.IsEmpty() {
		item.ItunesSummary = item.Summary
	}

	if item.GUID.IsEmpty() {
		if item.Enclosure != nil {
			item.GUID = NewGUID(item.Enclosure.URL)
		} else {
			item.GUID = NewGUID(item.File)
		}
	}

	// Extract information about file
	if f, err := os.Stat(item.File); !os.IsNotExist(err) {
		// file size
		item.FileSize = f.Size()
	}

}

// Validate channel
func (item *Item) Validate() error {
	log.Printf("Item[%s] Validate()...", item.Key)

	if item.File == "" {
		return fmt.Errorf("Item[%s] File path to audio file required", item.Key)
	}
	if _, err := os.Stat(item.File); os.IsNotExist(err) {
		return fmt.Errorf("Item[%s] %s", item.Key, err)
	}

	if item.PubDate.IsZero() {
		return fmt.Errorf("Item[%s] PubDate must be assigned", item.Key)
	}

	if item.Summary.IsEmpty() {
		return fmt.Errorf("Item[%s] Summary must be added", item.Key)
	}

	if item.Season > 0 && item.Episode == 0 {
		return fmt.Errorf("Item[%s] must have Episode if Season assigned", item.Key)
	}

	if item.Episode > 0 && item.Season == 0 {
		return fmt.Errorf("Item[%s] must have Season if Episode assigned", item.Key)
	}

	return nil
}
