package appcast

import (
	"encoding/xml"
)

type Sparkle struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XMLNSSparkle string   `xml:"xmlns:sparkle,attr"`
	XMLNSDC      string   `xml:"xmlns:dc,attr"`
	Channels     []SparkleChannel
}

type SparkleChannel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link,omitempty"`
	Description string   `xml:"description,omitempty"`
	Language    string   `xml:"language,omitempty"`
	Items       []SparkleItem
}

type SparkleItem struct {
	XMLName                 xml.Name            `xml:"item"`
	Title                   string              `xml:"title"`
	SparkleReleaseNotesLink string              `xml:"sparkle:releaseNotesLink,omitempty"`
	Description             *SparkleCdataString `xml:"description,omitempty"`
	PubDate                 string              `xml:"pubDate"`
	Enclosure               SparkleEnclosure    `xml:"enclosure,omitempty"`
}

// CdataString for XML CDATA
// See issue: https://github.com/golang/go/issues/16198
type SparkleCdataString struct {
	Value string `xml:",cdata"`
}

type SparkleItems []SparkleItem

type SparkleEnclosure struct {
	XMLName              xml.Name `xml:"enclosure"`
	URL                  string   `xml:"url,attr"`
	OS                   string   `xml:"sparkle:os,attr"`
	Version              string   `xml:"sparkle:version,attr"`
	DSASignature         string   `xml:"sparkle:dsaSignature,attr,omitempty"`
	EDSignature          string   `xml:"sparkle:edSignature,attr,omitempty"`
	InstallerArguments   string   `xml:"sparkle:installerArguments,attr,omitempty"`
	MinimumSystemVersion string   `xml:"sparkle:minimumSystemVersion,attr,omitempty"`
	Length               int      `xml:"length,attr,omitempty"`
	Type                 string   `xml:"type,attr"`
}
