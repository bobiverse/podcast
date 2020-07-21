package podcast

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Channel ..
type Channel struct {
	Domain string `xml:"-" yaml:"Domain"`

	// Text          string    `xml:",chardata" yaml:"-"`
	Link           string `xml:"link,omitempty" yaml:"Link"`
	Title          string `xml:"title" yaml:"Title"`
	Subtitle       string `xml:"subtitle,omitempty" yaml:"Subtitle"`
	Summary        *CDATA `xml:"summary,omitempty" yaml:"Summary"`
	Language       string `xml:"language,omitempty" yaml:"Language"`
	Author         string `xml:"author,omitempty" yaml:"Author"`
	Description    *CDATA `xml:"description,omitempty" yaml:"Description"`
	ContentEncoded *CDATA `xml:"content:encoded,omitempty" yaml:"ContentEncoded"`
	Image          *Image `xml:"image,omitempty" yaml:"Image"`

	// Docs about itunes https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	ItunesTitle    string    `xml:"itunes:title,omitempty" yaml:"ItunesTitle"`
	ItunesSubtitle string    `xml:"itunes:subtitle,omitempty" yaml:"ItunesSubtitle"`
	ItunesAuthor   string    `xml:"itunes:author,omitempty" yaml:"ItunesAuthor"`
	ItunesOwner    *Owner    `xml:"itunes:owner,omitempty" yaml:"Owner"`
	ItunesSummary  *CDATA    `xml:"itunes:summary,omitempty" yaml:"ItunesSummary"`
	ItunesType     string    `xml:"itunes:type,omitempty" yaml:"Type"`
	ItunesExplicit string    `xml:"itunes:explicit,omitempty" yaml:"Explicit"`
	ItunesKeywords string    `xml:"itunes:keywords,omitempty" yaml:"Keywords"`
	ItunesCategory *Category `xml:"itunes:category" yaml:"Category"`
	ItunesImage    *AttrHref `xml:"itunes:image" yaml:"ItunesImage"`

	LastBuildDate *Date  `xml:"lastBuildDate,omitempty" yaml:"LastBuildDate"`
	Copyright     string `xml:"copyright,omitempty" yaml:"Copyright"`

	Items ItemList `xml:"item" yaml:"Items"`
}

// Fix channel
func (channel *Channel) Fix() {

	if channel.Domain == "" && channel.Link != "" {
		u, err := url.Parse(channel.Link)
		if err == nil {
			channel.Domain = u.Scheme + "://" + u.Hostname()
		}
	}

	if channel.Link == "" && channel.Domain != "" {
		u, err := url.Parse(channel.Link)
		if err == nil {
			channel.Domain = u.Scheme + "://" + u.Hostname()
		}
	}

	// auto add last build time
	if channel.LastBuildDate == nil || channel.LastBuildDate.IsZero() {
		channel.LastBuildDate = &Date{time.Now()}
	}

	// Init as English podcast by default
	channel.Language = strings.ToLower(channel.Language)
	if channel.Language == "" {
		channel.Language = "en"
	}

	// Init as English podcast by default
	if channel.ContentEncoded.IsEmpty() {
		channel.ContentEncoded = channel.Description
	}

	if channel.Copyright == "" && channel.ItunesOwner != nil {
		channel.Copyright = fmt.Sprintf("℗ & © %s", channel.ItunesOwner.Name)
	}

	if channel.Image != nil && channel.Image.URL != "" {
		channel.Image.Title = channel.Title
		channel.Image.Link = channel.Link
	}

	// Copy generic fields to itunes
	if channel.ItunesTitle == "" {
		channel.ItunesTitle = channel.Title
	}
	if channel.ItunesSubtitle == "" {
		channel.ItunesSubtitle = channel.Subtitle
	}
	if channel.ItunesAuthor == "" {
		channel.ItunesAuthor = channel.Author
	}
	if channel.ItunesSummary.IsEmpty() {
		channel.ItunesSummary = channel.Summary
	}
	if channel.ItunesType == "" {
		channel.ItunesType = TypeEpisodic
	}
	if channel.ItunesExplicit == "" {
		channel.ItunesExplicit = ExplicitFalse
	}
	if channel.ItunesImage.IsEmpty() && !channel.Image.IsEmpty() {
		channel.ItunesImage = &AttrHref{channel.Image.URL}
	}

	// Fix items
	channel.Items.Fix(channel)

}

// Validate channel
func (channel *Channel) Validate() error {
	if !isValidURL(channel.Domain) {
		return fmt.Errorf("Invalid Domain. Please enter valid `Domain` or `Link` attribute")
	}

	if channel.Title == "" {
		return fmt.Errorf("Empty Channel Title")
	}

	if channel.Author == "" {
		return fmt.Errorf("Empty Channel Author")
	}

	// Apple Podcasts only supports values from the ISO 639 list (two-letter language codes, with some possible modifiers, such as "en-us").
	// https://www.loc.gov/standards/iso639-2/php/code_list.php
	// have dash - max length 5
	// no dash - max 2 letters
	// silly validator :)
	langLen := len(channel.Language)
	if langLen < 2 || langLen > 5 {
		return fmt.Errorf("Language must be in `ISO 639` format")
	} else if !strings.Contains(channel.Language, "-") && langLen != 2 {
		return fmt.Errorf("Language must be in `ISO 639` format")
	}

	if !isValidURL(channel.Link) {
		return fmt.Errorf("Empty Channel Link")
	}

	if channel.Summary.IsEmpty() {
		return fmt.Errorf("Empty Channel Summary")
	}

	if channel.Description.IsEmpty() {
		return fmt.Errorf("Empty Channel Description")
	}

	if channel.Image.IsEmpty() {
		return fmt.Errorf("Empty Channel Image URL")
	}

	linkURL, err := url.Parse(channel.Link)
	if err != nil {
		return fmt.Errorf("Error Channel Link (URL): %s", err)
	} else if linkURL.Scheme == "" || linkURL.Host == "" {
		return fmt.Errorf("Invalid Channel Link (URL) `%s`", linkURL)
	}

	if !inSlice(channel.ItunesType, TypeValues()) {
		return fmt.Errorf("Itunes Type must be one of the %v", TypeValues())
	}

	if !inSlice(channel.ItunesExplicit, ExplicitValues()) {
		return fmt.Errorf("Itunes Explicit must be one of the %v", ExplicitValues())
	}

	if channel.ItunesCategory.IsEmpty() {
		return fmt.Errorf("Empty Category. See: https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12")
	}

	if channel.ItunesOwner.IsEmpty() {
		return fmt.Errorf("Empty Owner. Add in format `My Name, my@email.xx`")
	}
	if !strings.Contains(channel.ItunesOwner.Email, "@") {
		return fmt.Errorf("Invalid Owner email field. Add Owner data in format `My Name, my@email.xx`")

	}

	if channel.Items.Len() == 0 {
		return fmt.Errorf("No episodes found. Add `Items:` into podcast yaml file")
	}

	return channel.Items.Validate()
}
