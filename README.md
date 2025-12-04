# greplab
Full code search in all your gitlab repositories

## Installation

```zsh
go install github.com/TWRaven/greplab@latest
```

Create a config.json file in $HOME/.greplab/ or /etc/greplab/

```json5
{
  // Required: Your GitLab Personal Access Token
  "gitlabToken": "glpat-xxxxxxxxxxxxxxxxx",
  // Optional: Gitlab API URL
  "gitlabUrl": "https://git.your-instance.com/api/v4",
}
```

All configuration options can be overridden with environment variables or cli flags.  
The priority is as follows:
- cli flag (highest prio)
- environment variable
- config key

| Cli flag       | Environment variable   | Config key (`json`) | Description                                                                       |
|:---------------|:-----------------------|:--------------------|:----------------------------------------------------------------------------------|
| `gitlab-token` | `GREPLAB_GITLAB_TOKEN` | `gitlabToken`       | Your Personal Access Token.                                                       |
| `gitlab-url`   | `GREPLAB_GITLAB_URL`   | `gitlabUrl`         | Base URL for self-hosted GitLab instances. (default: `https://gitlab.com/api/v4`) |

### Example
```zsh
GREPLAB_GITLAB_TOKEN="my-token" greplab search "func main"
```

## Usage
```zsh
greplab search <regex query> [flags]
```