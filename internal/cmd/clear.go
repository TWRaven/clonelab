package cmd

import (
	"fmt"
	"log/slog"
	"os/exec"
)

type Clear struct {
	Full bool `name:"full" short:"f" help:"Remove all files including cloned projects."`
}

func (c *Clear) Run(options *SharedOptions) error {
	if c.Full {
		slog.Info("removing clone directory", "directory", options.ProjectDir)

		cmd := exec.Command("rm", "-rf", options.ProjectDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not remove clone directory: %w", err)
		}

		return nil
	}

	slog.Info("removing cache file", "cache_file", options.GitlabCacheFile())

	cmd := exec.Command("rm", "-f", options.GitlabCacheFile())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not remove gitlab cache file: %w", err)
	}

	return nil
}
