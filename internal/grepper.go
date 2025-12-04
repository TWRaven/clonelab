package internal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProjectOutput struct {
	Project  *Project
	Output   string
	Err      error
	Duration time.Duration
}

func (project *ProjectOutput) PrettyPrint() {
	lines := strings.Split(project.Output, "\n")

	isNextLineFilename := true

	for _, line := range lines {
		// 1. If the line is empty, it means the next line will be a new file
		if line == "" {
			isNextLineFilename = true
			fmt.Println(line) // Print the empty line to keep spacing
			continue
		}

		// 2. If this is a filename, prefix it and reset the flag
		if isNextLineFilename {
			fmt.Printf("[%s] %s\n", project.Project.Name, line)
			isNextLineFilename = false
		} else {
			// 3. Otherwise, it's a code match line, print normally
			fmt.Println(line)
		}
	}
}

type Grepper struct {
	token            string
	projectsFilePath string
	numWorkers       int
}

func NewGrepper(gitlabToken string, projectsCacheDir string, numWorkers int) *Grepper {
	return &Grepper{
		token:            gitlabToken,
		projectsFilePath: projectsCacheDir,
		numWorkers:       numWorkers,
	}
}

type activeTracker struct {
	sync.Mutex
	jobs map[string]time.Time // project name to start time
}

func (g Grepper) Grep(projects []Project, pattern string) []ProjectOutput {
	bar := progressbar.NewOptions(len(projects),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Initializing..."),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)

	projectChan := make(chan Project, len(projects))
	outputChan := make(chan ProjectOutput, len(projects))
	var wg sync.WaitGroup

	tracker := &activeTracker{
		jobs: make(map[string]time.Time),
	}

	for range g.numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for project := range projectChan {
				tracker.Lock()
				tracker.jobs[project.Name] = time.Now()
				tracker.Unlock()

				start := time.Now()
				result := g.grepProject(project, pattern)

				tracker.Lock()
				delete(tracker.jobs, project.Name)
				tracker.Unlock()

				if result == nil {
					result = &ProjectOutput{}
				}
				result.Project = &project
				result.Duration = time.Since(start)
				outputChan <- *result
			}
		}()
	}

	for _, project := range projects {
		projectChan <- project
	}
	close(projectChan)

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	var finalResults []ProjectOutput
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	type slowJob struct {
		Name     string
		Duration time.Duration
	}

loop:
	for {
		select {
		case out, ok := <-outputChan:
			if !ok {
				break loop
			}
			bar.Add(1)
			if out.Err != nil || out.Output != "" {
				finalResults = append(finalResults, out)
			}

		case <-ticker.C:
			tracker.Lock()
			var slows []slowJob

			for name, startTime := range tracker.jobs {
				dur := time.Since(startTime)
				if dur > 2*time.Second {
					shortName := name
					if idx := strings.LastIndex(name, "/"); idx != -1 {
						shortName = name[idx+1:]
					}
					slows = append(slows, slowJob{Name: shortName, Duration: dur})
				}
			}
			tracker.Unlock()

			if len(slows) == 0 {
				bar.Describe("Ripping through repos...")
			} else {
				sort.Slice(slows, func(i, j int) bool {
					return slows[i].Duration > slows[j].Duration
				})

				var statusParts []string

				displayLimit := 5
				for i, job := range slows {
					if i >= displayLimit {
						statusParts = append(statusParts, fmt.Sprintf("+%d more", len(slows)-displayLimit))
						break
					}
					statusParts = append(statusParts, fmt.Sprintf("%s (%ds)", job.Name, int(job.Duration.Seconds())))
				}

				bar.Describe(fmt.Sprintf("Slow: %s", strings.Join(statusParts, ", ")))
			}
		}
	}

	return finalResults
}

func (g Grepper) clone(projectDir string, repoURL string) *ProjectOutput {
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return &ProjectOutput{Err: err}
	}

	cmd := exec.Command("git", "clone", repoURL, projectDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ProjectOutput{Err: fmt.Errorf("git clone failed: %s: %w", string(output), err)}
	}

	return nil
}

func (g Grepper) grepProject(project Project, pattern string) *ProjectOutput {
	repoURL := strings.Replace(project.URL, "https://", fmt.Sprintf("https://oauth2:%s@", g.token), 1)

	if err := os.MkdirAll(g.projectsFilePath, 0755); err != nil {
		return &ProjectOutput{Err: err}
	}

	splitRepoName := strings.Split(project.Name, "/")
	repoName := splitRepoName[len(splitRepoName)-1]

	projectDir := filepath.Join(g.projectsFilePath, repoName)

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		if errOutput := g.clone(projectDir, repoURL); errOutput != nil {
			return errOutput
		}
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = projectDir

		output, err := cmd.CombinedOutput()
		if err != nil {
			rmDirCmd := exec.Command("rm", "-rf", projectDir)
			_ = rmDirCmd.Run()

			if errOutput := g.clone(projectDir, repoURL); errOutput != nil {
				errOutput.Err = fmt.Errorf("git pull failed with error: %s: %w, tried to clone instead, which failed with error %w", string(output), err, errOutput)
				return errOutput
			}
		}
	}

	cmd := exec.Command("rg",
		"--pretty",
		"--sort=path",
		pattern,
	)
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Check if it's just "exit status 1" (no matches) vs a real error
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			// No matches found
			return nil
		}
		return &ProjectOutput{Err: fmt.Errorf("rg pull failed: %s: %w", string(output), err)}
	}

	return &ProjectOutput{
		Output: string(output),
	}
}
