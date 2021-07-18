package appcast

import (
	"bufio"
	"encoding/xml"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/russross/blackfriday/v2"
)

type Release struct {
	Name        string
	Description string
	Date        time.Time
	Version     string
	Prerelease  bool
	Assets      []*Asset
}

type Asset struct {
	URL  string
	Size int
}

type Options struct {
	InstallerArguments   string
	MinimumSystemVersion string
	RewriteURL           RewriteFunc
}

type Overrides map[OS]*Options

type AppCast struct {
	*Options
	Overrides Overrides

	AppcastURL      string
	Title           string
	Description     string
	WithPrereleases bool
	Source          Source
	IsMaxOS         MatchFunc
	IsWindows64     MatchFunc
	IsWindows32     MatchFunc
	Signatures      string
}

type Source interface {
	Releases() ([]*Release, error)
}

func (feed *AppCast) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var items SparkleItems

	releases, err := feed.Source.Releases()
	if err != nil {
		return err
	}

	sort.Slice(releases, func(i, j int) bool { return releases[i].Date.After(releases[j].Date) })

	for _, release := range releases {
		if release.Prerelease && !feed.WithPrereleases {
			log.Println("skip prelease", release.Version)
			continue
		}

		item, err := feed.releaseToSparkleItem(release)
		if err != nil {
			log.Println("warning:", err)
			continue
		}

		items = append(items, item...)
	}

	s := &Sparkle{
		Version:      "2.0",
		XMLNSSparkle: "http://www.andymatuschak.org/xml-namespaces/sparkle",
		XMLNSDC:      "http://purl.org/dc/elements/1.1/",
		Channels: []SparkleChannel{
			{
				Title:       feed.Title,
				Link:        feed.AppcastURL,
				Description: feed.Description,
				Items:       items,
			},
		},
	}

	return e.Encode(s)
}

var matchWhitespace = regexp.MustCompile(`\s+`)

func (feed *AppCast) releaseToSparkleItem(release *Release) (SparkleItems, error) {
	items := make([]SparkleItem, 0)
	var description *SparkleCdataString
	if release.Description != "" {
		htmlDescription := blackfriday.Run([]byte(release.Description))
		description = &SparkleCdataString{string(htmlDescription)}
	}

	var signatures string
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.URL, feed.Signatures) {
			signatures = asset.URL
		}

		os := feed.detectOS(asset.URL)
		if os == Unknown {
			continue
		}

		opt := new(Options)
		if feed.Overrides[os] != nil {
			*opt = *feed.Overrides[os]
		}
		if feed.Options != nil {
			_ = mergo.Merge(opt, feed.Options)
		}

		url := asset.URL
		if opt.RewriteURL != nil {
			url = opt.RewriteURL(url)
		}

		item := SparkleItem{
			Title:       release.Name,
			PubDate:     release.Date.UTC().Format(time.RFC1123),
			Description: description,
			Enclosure: SparkleEnclosure{
				Version:              strings.TrimPrefix(release.Version, "v"),
				URL:                  url,
				InstallerArguments:   opt.InstallerArguments,
				MinimumSystemVersion: opt.MinimumSystemVersion,
				Type:                 detectType(asset.URL),
				OS:                   os.String(),
				Length:               asset.Size,
			},
		}
		items = append(items, item)
	}

	if signatures != "" {
		r, err := http.Get(signatures)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		scanner := bufio.NewScanner(r.Body)
		for scanner.Scan() {
			sig := matchWhitespace.Split(scanner.Text(), -1)
			if len(sig) < 2 {
				continue
			}
			for i := range items {
				if strings.HasSuffix(items[i].Enclosure.URL, sig[0]) {
					items[i].Enclosure.DSASignature = sig[1]
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}

	}

	return items, nil
}

func (feed *AppCast) detectOS(url string) OS {
	if matchFallback(feed.IsMaxOS, isMacOS)(url) {
		return MacOS
	}
	if matchFallback(feed.IsWindows64, isWindows64)(url) {
		return Windows64
	}
	if matchFallback(feed.IsWindows32, isWindows32)(url) {
		return Windows32
	}

	return Unknown
}

func detectType(filename string) string {
	typ := mime.TypeByExtension(filepath.Ext(filename))
	if typ != "" {
		return typ
	}
	return "application/octet-stream"
}
