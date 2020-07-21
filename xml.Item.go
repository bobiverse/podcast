package podcast

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
)

// Item ..
type Item struct {
	Channel *Channel `xml:"-" yaml:"-"`
	Key     string   `xml:"-" yaml:"-"`

	// Text        string     `xml:",chardata"`
	Title       string     `xml:"title,omitempty" yaml:"Title"`
	Description *CDATA     `xml:"description,omitempty" yaml:"Description"`
	Encoded     *CDATA     `xml:"content:encoded,omitempty" yaml:"Encoded"`
	Summary     *CDATA     `xml:"summary,omitempty" yaml:"Summary"`
	Enclosure   *Enclosure `xml:"enclosure,omitempty" yaml:"Enclosure"`
	GUID        *GUID      `xml:"guid,omitempty" yaml:"GUID"`
	PubDate     *Date      `xml:"pubDate,omitempty" yaml:"PubDate"`
	Keywords    string     `xml:"itunes:keywords,omitempty" yaml:"Keywords"`
	Season      int        `xml:"itunes:season,omitempty" yaml:"Season"`
	Episode     int        `xml:"itunes:episode,omitempty" yaml:"Episode"`
	EpisodeType string     `xml:"itunes:episodeType,omitempty" yaml:"EpisodeType"`
	Explicit    string     `xml:"itunes:explicit,omitempty" yaml:"Explicit"`

	ItunesTitle   string `xml:"itunes:title,omitempty" yaml:"ItunesTitle"`
	ItunesSummary *CDATA `xml:"itunes:summary,omitempty" yaml:"ItunesSummary"`
	ItunesAuthor  string `xml:"itunes:author,omitempty" yaml:"Author"`

	// Different duration formats are accepted however it is recommended to convert the length of the episode into seconds.
	Duration Duration `xml:"itunes:duration,omitempty" yaml:"Duration"`

	File         string `xml:"-" yaml:"File"`
	FileSize     int64  `xml:"-" yaml:"FileSize"`
	FileMimeType string `xml:"-" yaml:"FileMimeType"`
	FileURL      string `xml:"-" yaml:"FileURL"`
}

// Weight of the item for sorting
// Seasons and episode taken
func (item *Item) Weight() int {
	weight := item.Season*1000 + item.Episode
	// log.Printf("WEIGHT: %d", weight)
	return weight
}

// ExtractKeyInfo info from key
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

	if item.File != "" {
		item.File = filepath.Clean(item.File)
	}

	if item.ItunesAuthor == "" {
		item.ItunesAuthor = item.Channel.ItunesAuthor
	}

	// Extract information about file
	if f, err := os.Stat(item.File); !os.IsNotExist(err) {
		// file size
		item.FileSize = f.Size()
		mime, err := mimetype.DetectFile(item.File)
		if err != nil {
			log.Printf("Warning: Couldn't get mime type of file `%s`. %s", item.File, err)
		} else {
			// fmt.Println(mime.String(), mime.Extension(), err)
			item.FileMimeType = mime.String()
		}
	}

	if item.FileURL == "" {
		item.FileURL = item.Channel.Domain + "/" + item.File

	}

	if item.Explicit == "" {
		item.Explicit = ExplicitFalse
	}

	item.Enclosure = &Enclosure{
		URL:    item.FileURL,
		Length: item.FileSize,
		Type:   item.FileMimeType,
	}

	if item.GUID.IsEmpty() {
		if item.Enclosure != nil {
			item.GUID = NewGUID(item.Enclosure.URL)
		} else {
			item.GUID = NewGUID(item.File)
		}
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

	if item.FileSize == 0 {
		return fmt.Errorf("Item[%s] FileSize required", item.Key)
	}

	if item.FileMimeType == "" {
		return fmt.Errorf("Item[%s] FileMimeType required", item.Key)
	}

	if item.PubDate.IsZero() {
		return fmt.Errorf("Item[%s] PubDate required", item.Key)
	}

	if item.Description.IsEmpty() {
		return fmt.Errorf("Item[%s] Description required", item.Key)
	}

	if item.Summary.IsEmpty() {
		return fmt.Errorf("Item[%s] Summary required", item.Key)
	}

	if item.Season > 0 && item.Episode == 0 {
		return fmt.Errorf("Item[%s] must have Episode if Season assigned", item.Key)
	}

	if item.Episode > 0 && item.Season == 0 {
		return fmt.Errorf("Item[%s] must have Season if Episode assigned", item.Key)
	}

	if !inSlice(item.Explicit, ExplicitValues()) {
		return fmt.Errorf("Item[%s] Explicit must be one of the %v", item.Key, ExplicitValues())
	}

	if item.Enclosure.IsEmpty() {
		return fmt.Errorf("Item[%s] Enclosure must be valid. Please input valid all of these: `File`, `FileSize`, `FileType`", item.Key)
	}

	if !isValidURL(item.Enclosure.URL) {
		return fmt.Errorf("Item[%s] Enclosure URL `%s` not valid. Please enter valid `FileURL`", item.Enclosure.URL, item.Key)
	}

	if item.Duration == 0 {
		return fmt.Errorf("Item[%s] Episode `Duration` required", item.Key)
	}

	return nil
}
