package internal

import (
	"fmt"
	"log/slog"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Project struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Topics    []string `json:"topics"`
	ShortName string   `json:"short_name"`

	URL      string `json:"url"`
	Archived bool   `json:"archived"`
}

func FetchProjects(c *gitlab.Client, programmingLanguage string) ([]Project, error) {
	slog.Info("Fetching projects from the gitlab api. This might take a while.")

	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	if programmingLanguage != "" {
		opt.WithProgrammingLanguage = &programmingLanguage
		slog.Debug("Filtering projects by programming language", "language", programmingLanguage)
	}

	var allProjects []Project
	for {
		projects, resp, err := c.Projects.ListProjects(opt)
		if err != nil {
			return nil, fmt.Errorf("error fetching projects: %v", err)
		}

		for _, p := range projects {
			allProjects = append(allProjects, Project{
				ID:        int(p.ID),
				Name:      p.PathWithNamespace,
				ShortName: p.Name,
				Topics:    p.Topics,
				Archived:  p.Archived,
				URL:       p.HTTPURLToRepo,
			})
		}

		// Exit if we're on the last page
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page
		opt.Page = resp.NextPage
	}

	return allProjects, nil
}
