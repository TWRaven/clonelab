package cmd

import (
	"log/slog"
	"os"

	"github.com/TWRaven/greplab/internal"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Search struct {
	Language string `help:"Programming language to filter projects on" short:"l"`

	GitlabToken string `help:"Gitlab access token" short:"t"`
	GitlabURL   string `help:"Gitlab URL" default:"https://gitlab.com/api/v4" short:"u"`
	NumWorkers  int    `help:"Number of workers" default:"5" short:"w"`

	ExcludeProjects []string `help:"Exclude projects" short:"x"`

	Query string `arg:"" name:"query" help:"Regex search query"`
}

func (s *Search) Run(options *SharedOptions) error {
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

	grepper := internal.NewGrepper(s.GitlabToken, options.ProjectsCacheDir(), s.NumWorkers)
	output := grepper.Grep(projects, s.Query)

	for _, project := range output {
		if project.Err != nil {
			slog.Error("Failed to grep project ", project.Project.Name, " with error: ", project.Err)
			continue
		}

		project.PrettyPrint()
	}

	return nil
}
