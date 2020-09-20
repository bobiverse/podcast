package podcast

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
)

// Item ..
type Item struct {
	Channel *Channel `xml:"-" yaml:"-"`
	Key     string   `xml:"-" yaml:"-"`

	// Text        string     `xml:",chardata"`
	Title string `xml:"title,omitempty" yaml:"Title"`

	// A single, descriptive sentence for your podcast or episode in <itunes:subtitle>.
	Subtitle string `xml:"itunes:subtitle,omitempty" yaml:"Subtitle"`

	// One or more sentences, or a paragraph, describing your podcast or episode in <description>.
	// Apple recommends the text in <description> be the same as the text in <content:encoded>, but in plain text form.
	Description *CDATA `xml:"description,omitempty" yaml:"Description"`

	// Apple recommends the text in <content:encoded> be the same as the text in <description>, but in HTML.
	ContentEncoded *CDATA `xml:"content:encoded,omitempty" yaml:"Encoded"`

	Enclosure   *Enclosure `xml:"enclosure,omitempty" yaml:"Enclosure"`
	Link        string     `xml:"link,omitempty" yaml:"Link"`
	GUID        *GUID      `xml:"guid,omitempty" yaml:"GUID"`
	PubDate     *Date      `xml:"pubDate,omitempty" yaml:"PubDate"`
	Keywords    string     `xml:"itunes:keywords,omitempty" yaml:"Keywords"`
	Season      int        `xml:"itunes:season,omitempty" yaml:"Season"`
	Episode     int        `xml:"itunes:episode,omitempty" yaml:"Episode"`
	EpisodeType string     `xml:"itunes:episodeType,omitempty" yaml:"EpisodeType"`
	Explicit    string     `xml:"itunes:explicit,omitempty" yaml:"Explicit"`

	// One or more sentences summarizing your podcast or episode in <itunes:summary>.
	ItunesSummary *CDATA    `xml:"itunes:summary,omitempty" yaml:"Summary"`
	ItunesAuthor  string    `xml:"itunes:author,omitempty" yaml:"Author"`
	ItunesImage   *AttrHref `xml:"itunes:image,omitempty" yaml:"Image"`

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

	episode, err := strconv.Atoi(item.Key[4:6])
	if err != nil {
		log.Printf("Warning: Episode can't be extracted from Item [%s]", item.Key)

	} else if item.Episode == 0 {
		// assign if not assigned in yaml file
		item.Episode = episode
	}
}

// Fix ..
func (item *Item) Fix() {
	// log.Printf("Item[%s] Fix()...", item.Key)

	if item.ContentEncoded.IsEmpty() && !item.Description.IsEmpty() && item.Description != item.ContentEncoded {
		item.ContentEncoded = item.Description
	}

	if item.ContentEncoded == item.Description {
		item.ContentEncoded = &CDATA{Text: "<p>" + item.Description.String() + "</p>"}
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
		item.FileURL = pathToURL(item.Channel.Domain, item.File)
	}
	if item.FileURL != "" && !isValidURL(item.FileURL) {
		item.FileURL = pathToURL(item.Channel.Domain, item.FileURL)
	}

	if item.Explicit == "" {
		item.Explicit = ExplicitYes
	}

	if item.EpisodeType == "" {
		item.EpisodeType = EpisodeTypeFull
	}

	if item.ItunesImage.IsEmpty() {
		item.ItunesImage = item.Channel.ItunesImage
	} else if !isValidURL(item.ItunesImage.Href) {
		item.ItunesImage.Href = pathToURL(item.Channel.Domain, item.ItunesImage.Href)
	}

	// Try detect duration automatically
	if item.Duration == 0 && item.File != "" {
		var buf []byte
		var re *regexp.Regexp
		var matches [][]byte

		// ffprobe
		_, buf, _ = runBash("ffprobe", item.File)
		re = regexp.MustCompile("Duration: ([0-9:]+)")
		matches = re.FindSubmatch(buf)
		if len(matches) >= 2 {
			item.Duration.Set(string(matches[1]))
		}

		// fmpeg
		if item.Duration == 0 {
			_, buf, _ = runBash("ffmpeg", "-i", item.File, "2>&1")
			re = regexp.MustCompile("Duration: ([0-9:]+)")
			matches = re.FindSubmatch(buf)
			if len(matches) >= 2 {
				item.Duration.Set(string(matches[1]))
			}
		}

		// exiftool
		if item.Duration == 0 {
			buf, _, _ = runBash("exiftool", item.File)
			re = regexp.MustCompile("Duration.+?: ([0-9:]+)")
			matches = re.FindSubmatch(buf)
			if len(matches) >= 2 {
				item.Duration.Set(string(matches[1]))
			}
		}

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

	if item.Link == "" {
		item.Link = item.Enclosure.URL
	}

}

// Validate channel
func (item *Item) Validate() error {
	// log.Printf("Item[%s] Validate()...", item.Key)

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

	if item.Link == "" {
		return fmt.Errorf("Item[%s] `Link` required", item.Key)
	}

	if item.GUID.IsEmpty() {
		return fmt.Errorf("Item[%s] `GUID` required", item.Key)
	}

	// if item.Summary.IsEmpty() {
	// 	return fmt.Errorf("Item[%s] Summary required", item.Key)
	// }

	if item.Season > 0 && item.Episode == 0 {
		return fmt.Errorf("Item[%s] must have Episode if Season assigned", item.Key)
	}

	if item.Episode > 0 && item.Season == 0 {
		return fmt.Errorf("Item[%s] must have Season if Episode assigned", item.Key)
	}

	if !inSlice(item.Explicit, ExplicitValues()) {
		return fmt.Errorf("Item[%s] Explicit must be one of the %v", item.Key, ExplicitValues())
	}

	if !inSlice(item.EpisodeType, EpisodeTypesValues()) {
		return fmt.Errorf("Item[%s] EpisodeType must be one of the %v", item.Key, EpisodeTypesValues())
	}

	if item.Enclosure.IsEmpty() {
		return fmt.Errorf("Item[%s] Enclosure must be valid. Please input valid all of these: `File`, `FileSize`, `FileType`", item.Key)
	}

	if !isValidURL(item.Enclosure.URL) {
		return fmt.Errorf("Item[%s] Enclosure URL `%s` not valid. Please enter valid `FileURL`", item.Enclosure.URL, item.Key)
	}

	if item.Duration == 0 {
		return fmt.Errorf("Item[%s] Episode `Duration` required. Add it manualy or install `ffprobe`, `fmpeg` or `exiftool` to get duration automatically", item.Key)
	}

	if item.ItunesImage != nil && !isValidURL(item.ItunesImage.Href) {
		return fmt.Errorf("Item[%s] Episode `Image` must be valid URL", item.Key)
	}

	return nil
}
