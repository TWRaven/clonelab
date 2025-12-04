package cmd

import "path/filepath"

type SharedOptions struct {
	GreplabDir string
}

func (s SharedOptions) GitlabCacheFile() string {
	return filepath.Join(s.GreplabDir, "cache.json")
}

func (s SharedOptions) ProjectsCacheDir() string {
	return filepath.Join(s.GreplabDir, "projects")
}
