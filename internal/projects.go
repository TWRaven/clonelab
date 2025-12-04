package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type ProjectClient struct {
	client *gitlab.Client

	cacheFilePath       string
	programmingLanguage string
	excludeProjects     map[string]struct{}
}

func NewProjectClient(client *gitlab.Client, gitlabCacheFile string, language string, excludeProjects []string) *ProjectClient {
	excluded := make(map[string]struct{})
	for _, project := range excludeProjects {
		excluded[project] = struct{}{}
	}

	return &ProjectClient{
		client:              client,
		cacheFilePath:       gitlabCacheFile,
		programmingLanguage: language,
		excludeProjects:     excluded,
	}
}

func (p *ProjectClient) GetProjects() ([]Project, error) {
	var projects []Project

	if _, err := os.Stat(p.cacheFilePath); err == nil {
		file, err := os.Open(p.cacheFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&projects)
		if err != nil {
			return nil, err
		}
	} else {
		projects, err := FetchProjects(p.client, p.programmingLanguage)
		if err != nil {
			return nil, err
		}

		// create cache file if it doesn't yet exist
		dir := filepath.Dir(p.cacheFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		file, err := os.Create(p.cacheFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		err = encoder.Encode(projects)
		if err != nil {
			return nil, err
		}
	}

	slices.DeleteFunc(projects, func(item Project) bool {
		_, ok := p.excludeProjects[item.ShortName]
		return ok
	})
	return projects, nil
}
