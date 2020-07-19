package podcast

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

// PodcastYAMLFilename ..
const PodcastYAMLFilename = "podcast.yml"

// Podcast ..
type Podcast struct {
	ContentPath string `yaml:"-"`
	OutputPath  string `yaml:"-"`

	Title string `yaml:"Title"`

	Feed *XMLRoot `yaml:"-"`
}

// New ..
func New(contentPath string) (*Podcast, error) {
	log.Println("[podcast] New ")

	podcast := &Podcast{
		Feed: &XMLRoot{
			Itunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
			Content: "http://purl.org/rss/1.0/modules/content/",
			Atom:    "http://www.w3.org/2005/Atom",
			Version: "2.0",
			Channel: &Channel{},
		},
	}

	// load content from given directory
	if contentPath != "" {
		if err := podcast.LoadFromPath(contentPath); err != nil {
			return podcast, err
		}
	}

	return podcast, nil
}

// LoadFromPath ..
func (podcast *Podcast) LoadFromPath(contentPath string) error {

	// check if folder exists
	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		return err
	}
	podcast.ContentPath = contentPath

	// check if podcast YAML file exists
	fpath := filepath.Join(contentPath, PodcastYAMLFilename)
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return err
	}

	// Read YAML contents
	buf, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	// Parse YAML into struct
	if err := yaml.Unmarshal(buf, &podcast.Feed.Channel); err != nil {
		return err
	}

	// podcast.Feed.Channel.Title = podcast.Title

	color.Magenta("%+v", podcast.Feed.Channel)
	return nil
}

// Fix misconfigs and populate empty values with defaults  before saving ..
func (podcast *Podcast) Fix() {
	log.Println("[podcast] Fix ")
	podcast.Feed.Channel.Fix()
}

// Validate before saving ..
func (podcast *Podcast) Validate() error {
	log.Println("[podcast] Validate ")

	if err := podcast.Feed.Channel.Validate(); err != nil {
		return err
	}

	return nil
}

// SaveToFile ..
func (podcast *Podcast) SaveToFile() error {
	// fix some values
	podcast.Fix()

	// validate feed before saving to file
	if err := podcast.Validate(); err != nil {
		return err
	}

	// generate XML and save to file
	log.Println("[podcast] SaveToFile ")

	buf, err := podcast.Feed.ToXML("")
	if err != nil {
		return err
	}

	color.Yellow("%s", buf)
	return nil
}
