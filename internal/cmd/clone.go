package cmd

import (
	"log/slog"
	"os"

	"github.com/TWRaven/clonelab/internal"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Clone struct {
	Language string `help:"Programming language to filter projects on" short:"l"`

	GitlabToken string `help:"Gitlab access token" short:"t"`
	GitlabURL   string `help:"Gitlab URL" default:"https://gitlab.com/api/v4" short:"u"`
	NumWorkers  int    `help:"Number of workers" default:"5" short:"w"`

	ExcludeProjects []string `help:"Exclude projects" short:"x"`
}

func (s *Clone) Run(options *SharedOptions) error {
	gitlabClient, err := gitlab.NewClient(s.GitlabToken, gitlab.WithBaseURL(s.GitlabURL))
	if err != nil {
		slog.Error("Failed to create GitLab client with error: ", err)
		os.Exit(1)
	}

	projectsClient := internal.NewProjectClient(gitlabClient, options.GitlabCacheFile(), s.Language, s.ExcludeProjects)
	projects, err := projectsClient.GetProjects()

	if err != nil {
		slog.Error("Failed to get projects with error: ", err)
		os.Exit(1)
	}

	cloner := internal.NewCloner(s.GitlabToken, options.ProjectsDir(), s.NumWorkers)
	output := cloner.Clone(projects)

	for _, project := range output {
		if project.Err != nil {
			slog.Error("Failed to clone project", "project", project.Project.Name, "error", project.Err)
		}
	}

	slog.Info("Finished cloning projects", "directory", options.ProjectsDir())

	return nil
}
