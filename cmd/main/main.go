package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/source/gitlab"
)

func main() {
	feed := &appcast.AppCast{
		Title:       "Example",
		Description: "Updates for example app.",
		AppcastURL:  "https://example.com/appcast.xml",
		Overrides: appcast.Overrides{
			appcast.Windows64: {InstallerArguments: "/passive"},
			appcast.Windows32: {InstallerArguments: "/passive"},
			appcast.MacOS:     {MinimumSystemVersion: "10.13.0"},
		},
		Options: &appcast.Options{
			RewriteURL: appcast.RegexRewrite("^.*/([^/]+)$", "https://dl.example.com/$1"),
		},
		IsMaxOS:     appcast.RegexMatch(`^.*\.dmg$`),
		IsWindows64: appcast.GlobMatch(`*64-bit.msi`),
		IsWindows32: appcast.GlobMatch(`**/*32-bit*.msi`),
		Signatures:  "signatures.txt",
		Source: gitlab.New(&gitlab.Config{
			Repo:  "org/repo",
			Token: os.Getenv("GITLAB_TOKEN"),
		}),
	}

	buf := bytes.NewBuffer(nil)
	buf.Write([]byte(xml.Header))

	xw := xml.NewEncoder(buf)
	xw.Indent("", "  ")

	err := xw.Encode(feed)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("appcast.xml", buf.Bytes(), 0755)
	if err != nil {
		log.Fatal(err)
	}
}
