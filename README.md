# clonelab
Clone all your gitlab repositories

## Installation

```zsh
go install github.com/TWRaven/clonelab@latest
```

Create a config.json file in $HOME/.clonelab/ or /etc/clonelab/

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
| `gitlab-token` | `CLONELAB_GITLAB_TOKEN` | `gitlabToken`       | Your Personal Access Token.                                                       |
| `gitlab-url`   | `CLONELAB_GITLAB_URL`   | `gitlabUrl`         | Base URL for self-hosted GitLab instances. (default: `https://gitlab.com/api/v4`) |

### Example
```zsh
CLONELAB_GITLAB_TOKEN="my-token" clonelab clone
```

## Usage
```zsh
clonelab clone [flags]
```