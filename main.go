package main

import (
	"log/slog"
	"os"

	"github.com/TWRaven/clonelab/internal/cmd"
	"github.com/alecthomas/kong"
)

type CLI struct {
	Dir       string `type:"path" help:"clonelab directory" default:"/tmp/clonelab" short:"d"`
	Verbosity string `help:"Log verbosity" short:"v" default:"info" enum:"debug,info,warn,error"`

	Clone cmd.Clone `cmd:"" help:"Clone occurrences in gitlab projects."`
	Clear cmd.Clear `cmd:"" help:"Clear gitlab cache files."`
}

func main() {
	cli := CLI{}
	kongCtx := kong.Parse(
		&cli,
		kong.Name("clonelab"),
		kong.Description("A GitLab clone tool"),
		kong.DefaultEnvars("clonelab"),
		kong.Configuration(kong.JSON, "~/.clonelab/config.json", "/etc/clonelab/config.json"),
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
		ProjectDir: cli.Dir,
	})
	kongCtx.FatalIfErrorf(err)
}
