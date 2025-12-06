package cmd

import "path/filepath"

type SharedOptions struct {
	ProjectDir string
}

func (s SharedOptions) GitlabCacheFile() string {
	return filepath.Join(s.ProjectDir, "cache.json")
}

func (s SharedOptions) ProjectsDir() string {
	return filepath.Join(s.ProjectDir, "projects")
}
