package podcast

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

// Channel ..
type Channel struct {
	Domain string `xml:"-" yaml:"Domain"`

	SelfLink *AttrHref `xml:"atom:link,omitempty" yaml:"SelfLink"`

	// Text          string    `xml:",chardata" yaml:"-"`
	Link           string `xml:"link,omitempty" yaml:"Link"`
	Title          string `xml:"title" yaml:"Title"`
	Subtitle       string `xml:"itunes:subtitle,omitempty" yaml:"Subtitle"`
	Language       string `xml:"language,omitempty" yaml:"Language"`
	Description    *CDATA `xml:"description,omitempty" yaml:"Description"`
	ContentEncoded *CDATA `xml:"content:encoded,omitempty" yaml:"ContentEncoded"`
	Image          *Image `xml:"image,omitempty" yaml:"Image"`

	// Docs about itunes https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	ItunesTitle    string    `xml:"itunes:title,omitempty" yaml:"ItunesTitle"`
	ItunesAuthor   string    `xml:"itunes:author,omitempty" yaml:"Author"`
	ItunesOwner    *Owner    `xml:"itunes:owner,omitempty" yaml:"Owner"`
	ItunesSummary  *CDATA    `xml:"itunes:summary,omitempty" yaml:"Summary"`
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

	// Try to get `Domain` from `Link`
	if channel.Domain == "" && channel.Link != "" {
		u, err := url.Parse(channel.Link)
		if err == nil {
			channel.Domain = u.Scheme + "://" + u.Hostname()
		}
	}

	// Fix `Domain` if not valid URL
	// `example.org/` ==> `https://example.org` (yes, using SSL)
	if channel.Domain != "" && !isValidURL(channel.Domain) {
		s := "https://" + channel.Domain
		u, err := url.Parse(s)
		if err == nil {
			channel.Domain = u.Scheme + "://" + u.Hostname()
		}
	}

	// Try to get `Link` from `Domain`
	if channel.Link == "" && channel.Domain != "" {
		u, err := url.Parse(channel.Domain)
		if err == nil {
			channel.Link = u.Scheme + "://" + u.Hostname()
		}
	}

	// Fix `Link` if not valid URL
	// `./my-link` ==> `https://example.org/my-link`
	if channel.Link != "" && !isValidURL(channel.Link) {
		channel.Link = pathToURL(channel.Domain, channel.Link)
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

		if !isValidURL(channel.Image.URL) {
			channel.Image.URL = pathToURL(channel.Domain, channel.Image.URL)
		}
	}

	// Copy generic fields to itunes
	if channel.ItunesTitle == "" {
		channel.ItunesTitle = channel.Title
	}
	if channel.ItunesType == "" {
		channel.ItunesType = PodcastTypeEpisodic
	}
	if channel.ItunesExplicit == "" {
		channel.ItunesExplicit = ExplicitFalse
	}
	if channel.ItunesImage.IsEmpty() && !channel.Image.IsEmpty() {
		channel.ItunesImage = &AttrHref{Href: channel.Image.URL}
	}

	if !channel.SelfLink.IsEmpty() && !isValidURL(channel.SelfLink.Href) {
		channel.SelfLink.Href = pathToURL(channel.Domain, channel.SelfLink.Href)
		channel.SelfLink.Rel = "self"
		channel.SelfLink.Type = "application/rss+xml"
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

	if channel.ItunesAuthor == "" {
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

	if channel.ItunesSummary.IsEmpty() {
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

	if !inSlice(channel.ItunesType, PodcastTypeValues()) {
		return fmt.Errorf("Itunes Type must be one of the %v", PodcastTypeValues())
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

	if channel.SelfLink.IsEmpty() {
		log.Printf("Warning: It's recommended to add this RSS feed URL in `SelfLink` param ")
	}

	if channel.Items.Len() == 0 {
		return fmt.Errorf("No episodes found. Add `Items:` into podcast yaml file")
	}

	return channel.Items.Validate()
}
