package bintray

import "time"

type Package struct {
	Desc     string   `json:"desc"`
	Versions []string `json:"versions"`
}

type PackageFile struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Version string    `json:"version"`
	Created time.Time `json:"created"`
	Size    int       `json:"size"`
}

type Version struct {
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Published bool      `json:"published"`
	Created   time.Time `json:"created"`
}

type VersionFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
}
