package main

import (
	"log/slog"
	"os"

	"github.com/TWRaven/greplab/internal/cmd"
	"github.com/alecthomas/kong"
)

type CLI struct {
	CacheDir  string `type:"path" help:"Greplab cache directory" default:"/tmp/greplab" short:"c"`
	Verbosity string `help:"Log verbosity" short:"v" default:"info" enum:"debug,info,warn,error"`

	Search cmd.Search `cmd:"" help:"Search occurrences in gitlab projects."`
	Clear  cmd.Clear  `cmd:"" help:"Clear gitlab cache files."`
}

func main() {
	cli := CLI{}
	kongCtx := kong.Parse(
		&cli,
		kong.Name("greplab"),
		kong.Description("A GitLab Ripgrep tool"),
		kong.DefaultEnvars("GREPLAB"),
		kong.Configuration(kong.JSON, "~/.greplab/config.json", "/etc/greplab/config.json"),
	)

	var level slog.Level
	if err := level.UnmarshalText([]byte(cli.Verbosity)); err != nil {
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	err := kongCtx.Run(&cmd.SharedOptions{
		GreplabDir: cli.CacheDir,
	})
	kongCtx.FatalIfErrorf(err)
}
