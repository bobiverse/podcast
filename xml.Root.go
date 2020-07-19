package podcast

import (
	"encoding/xml"
)

// XMLFilePrefix - feed top line
const XMLFilePrefix = `<?xml version="1.0" encoding="UTF-8"?>`

// XMLRoot - rss feed base
type XMLRoot struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Itunes  string   `xml:"xmlns:itunes,attr,omitempty"`
	Content string   `xml:"xmlns:content,attr,omitempty"`
	Atom    string   `xml:"xmlns:atom,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
	Channel *Channel `xml:"channel"`
}

// ToXML ..
func (feed *XMLRoot) ToXML(xmlPrefix string) ([]byte, error) {
	if xmlPrefix == "" {
		xmlPrefix = XMLFilePrefix
	}

	buf, err := xml.MarshalIndent(feed, "  ", "    ")
	if err != nil {
		return nil, err
	}

	buf = append([]byte(xmlPrefix+"\n"), buf...)

	return buf, nil
}
