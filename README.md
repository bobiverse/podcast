# podcast
Generate valid podcast _XML_ from given _YAML_ file with minimum details.

## Validates against
- https://castfeedvalidator.com/
- https://podba.se/validate/

## Accepted by
- Apple Itunes
- Google podcasts
- Spotify

## Auto-populated fields
Large amount of fields are pre-filled or auto-populated. If you see bug or you need different value in this field you can overwrite it in your _YAML_ file.

For example if you want to change default value of `<itunes:explicit>false</itunes:explicit>` use _CamelCased_ field name in your _YAML_ file as:
```yaml
Explicit: true
```

To see what are YAML field names for XML tags see code  `podcast/xml.Item.go` and `podcast/xml.Channel.go`.

---
```go
// variable: `ItunesImage`
// YAML: `Image:`
// XML: `<itunes:image>`
ItunesImage   *AttrHref `xml:"itunes:image,omitempty" yaml:"Image"`
```

## Example
```go
package main

import (
	"log"

	"github.com/briiC/podcast"
)

func main() {
	// Load YML file required to generate XML
	Podcast, err := podcast.New("./podcast.yml")
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

	// Get XML to output
	buf, _ := Podcast.XML()
	log.Printf("XML:\n%s", Podcast.XML())

	// -- OR -- save XML to file
	if err := Podcast.SaveToFile("feed.xml"); err != nil {
		log.Printf("ERROR: %s", err)
	}

	// List episodes in your code or send to HTML template
	for _, ep := range Podcast.Episodes() {
		log.Printf(">> %s", ep.Title)
	}
}

```
and

```yaml
Title: Podcast example
Domain: https://exampple.xx
# Link: https://exampple.xx/my/podcast/
Author: Neo and Trinity
Owner: John, john@example.xx
Description: Long description of this podcast. Couple of sentences. Or more.
Summary: Very short description of this podcast
Language: en

# Absolute URL or Relative (domain will be prepended)
Image: /podcast.png

# Category. See https://help.apple.com/itc/podcasts_connect/#/itc9267a2f12
# Comma separated. Use first as primary cateogory.
Category: TV & Film, TV Reviews

# Subtitle: This is just an example
# Keywords: tv, reviews ,example

Items:
    S01E02:
        Title: John Wick - Chapter 3 - Parabellum
        Description: Reviewing movie "John Wick - Chapter 3 - Parabellum" (7.5/10)
        File: ./episodes/S01E02.mp3
        PubDate: 2020-07-14
    S01E01:
        Title: Apocalypse - The Second World War
        File: ./episodes/S01E01.mp3
        FileURL: https://exampple.xx/different/path/to/public/file/S01E01.mp3
        PubDate: 2020-07-07
        Description: Reviewing movie "Apocalypse - The Second World War" (9/10)
        Image: ./images/custom-episode-image.jpeg

```
generates _XML_

```xml
<?xml version="1.0" encoding="UTF-8"?>
  <rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
	  <lastBuildDate>Thu, 23 Jul 2020 09:36:30 EEST</lastBuildDate>
      <channel>
          <link>https://exampple.xx</link>
          <title>Podcast example</title>
          <language>en</language>
          <description><![CDATA[Long description of this podcast. Couple of sentences. Or more.]]></description>
          <content:encoded><![CDATA[Long description of this podcast. Couple of sentences. Or more.]]></content:encoded>
          <image>
              <url>https://exampple.xx/podcast.png</url>
              <title>Podcast example</title>
              <link>https://exampple.xx</link>
          </image>
          <itunes:title>Podcast example</itunes:title>
          <itunes:author>Neo and Trinity</itunes:author>
          <itunes:owner>
              <itunes:name>John</itunes:name>
              <itunes:email>john@example.xx</itunes:email>
          </itunes:owner>
          <itunes:summary><![CDATA[Very short description of this podcast]]></itunes:summary>
          <itunes:type>episodic</itunes:type>
          <itunes:explicit>false</itunes:explicit>
          <itunes:category text="TV &amp; Film">
              <itunes:category text="TV Reviews"></itunes:category>
          </itunes:category>
          <itunes:image href="https://exampple.xx/podcast.png"></itunes:image>
          <lastBuildDate>Thu, 23 Jul 2020 00:08:47 EEST</lastBuildDate>
          <copyright>℗ &amp; © John</copyright>
          <item>
              <title>Apocalypse - The Second World War</title>
              <description><![CDATA[Reviewing movie "Apocalypse - The Second World War" (9/10)]]></description>
              <content:encoded><![CDATA[<p>Reviewing movie "Apocalypse - The Second World War" (9/10)</p>]]></content:encoded>
              <enclosure url="https://exampple.xx/different/path/to/public/file/S01E01.mp3" length="72552696" type="audio/mpeg"></enclosure>
              <link>https://exampple.xx/different/path/to/public/file/S01E01.mp3</link>
              <guid isPermaLink="true">https://exampple.xx/different/path/to/public/file/S01E01.mp3</guid>
              <pubDate>Tue, 07 Jul 2020 00:00:00 UTC</pubDate>
              <itunes:season>1</itunes:season>
              <itunes:episode>1</itunes:episode>
              <itunes:episodeType>full</itunes:episodeType>
              <itunes:explicit>false</itunes:explicit>
              <itunes:author>Neo and Trinity</itunes:author>
              <itunes:image href="https://exampple.xx/images/custom-episode-image.jpeg"></itunes:image>
              <itunes:duration>3022</itunes:duration>
          </item>
          <item>
              <title>John Wick - Chapter 3 - Parabellum</title>
              <description><![CDATA[Reviewing movie "John Wick - Chapter 3 - Parabellum" (7.5/10)]]></description>
              <content:encoded><![CDATA[<p>Reviewing movie "John Wick - Chapter 3 - Parabellum" (7.5/10)</p>]]></content:encoded>
              <enclosure url="https://exampple.xx/episodes/S01E02.mp3" length="78107374" type="audio/mpeg"></enclosure>
              <link>https://exampple.xx/episodes/S01E02.mp3</link>
              <guid isPermaLink="true">https://exampple.xx/episodes/S01E02.mp3</guid>
              <pubDate>Tue, 14 Jul 2020 00:00:00 UTC</pubDate>
              <itunes:season>1</itunes:season>
              <itunes:episode>2</itunes:episode>
              <itunes:episodeType>full</itunes:episodeType>
              <itunes:explicit>false</itunes:explicit>
              <itunes:author>Neo and Trinity</itunes:author>
              <itunes:image href="https://exampple.xx/podcast.png"></itunes:image>
              <itunes:duration>3254</itunes:duration>
          </item>
      </channel>
  </rss>
```


## TODO:
- tests
- keep clean XML (?) remove tags with default values already
- calculate audio file duration by sample rate and sample count
- parse `PubDate` from different datetime formats
