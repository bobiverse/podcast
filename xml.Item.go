package podcast

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Item ..
type Item struct {
	Key  string `xml:"-" yaml:"-"`
	File string `xml:"-" yaml:"File"`

	// Text        string     `xml:",chardata"`
	Title       string     `xml:"title,omitempty" yaml:"Title"`
	Description string     `xml:"description,omitempty" yaml:"Description"`
	Encoded     string     `xml:"encoded,omitempty" yaml:"Encoded"`
	Author      string     `xml:"author,omitempty" yaml:"Author"`
	Summary     string     `xml:"summary,omitempty" yaml:"Summary"`
	Enclosure   *Enclosure `xml:"enclosure,omitempty" yaml:"Enclosure"`
	GUID        *GUID      `xml:"guid,omitempty" yaml:"GUID"`
	PubDate     string     `xml:"pubDate,omitempty" yaml:"PubDate"`
	Duration    string     `xml:"duration,omitempty" yaml:"Duration"`
	Keywords    string     `xml:"keywords,omitempty" yaml:"Keywords"`
	Season      int        `xml:"season,omitempty" yaml:"Season"`
	Episode     int        `xml:"episode,omitempty" yaml:"Episode"`
	EpisodeType string     `xml:"episodeType,omitempty" yaml:"EpisodeType"`
	Explicit    string     `xml:"explicit,omitempty" yaml:"Explicit"`
}

// Weight of the item for sorting
// Seasons and episode taken
func (item *Item) Weight() int {
	return item.Season*1000 + item.Episode
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

	if item.Summary == "" {
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
