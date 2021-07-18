package gitlab

import (
	"github.com/abemedia/appcast"
	"github.com/xanzy/go-gitlab"
)

type Config struct {
	Token string
	Repo  string
}

type gitlabSource struct {
	client *gitlab.Client
	repo   string
}

func New(c *Config) appcast.Source {
	s := new(gitlabSource)

	git, err := gitlab.NewClient(c.Token)
	if err != nil {
		panic(err)
	}

	s.client = git
	s.repo = c.Repo

	return s
}

func (s *gitlabSource) Releases() ([]*appcast.Release, error) {
	pkg, _, err := s.client.Releases.ListReleases(s.repo, nil)
	if err != nil {
		return nil, err
	}

	var result []*appcast.Release
	for _, file := range pkg {
		assets := make([]*appcast.Asset, len(file.Assets.Links))
		for i, l := range file.Assets.Links {
			assets[i] = &appcast.Asset{URL: l.URL}
		}
		r := &appcast.Release{
			Name:        file.Name,
			Description: file.Description,
			Version:     file.TagName,
			Date:        *file.CreatedAt,
			Assets:      assets,
		}
		result = append(result, r)
	}

	return result, nil
}
