package main

import (
	"os"
	"github.com/alecthomas/kong"
)

type InputMode uint8
const (
	InputScript InputMode = iota
	InputInline
)

type Action uint8
const (
	ActionRun Action = iota
	ActionBuild
	ActionMigrate
	ActionExplain
	ActionClearCache
)

type Config struct {
	Input  InputMode
	Action Action

	ScriptPath string
	Args       []string
	InlineCode string

	OutputPath string
	MigrateDir string
	FieldSep   string
	Parallel   int
	NoCache    bool
}

type cliArgs struct {
	FieldSep   string `short:"F" help:"Field separator for f/fields in pipe mode (default: whitespace split)."`
	Out        string `short:"o" help:"Build output path."`
	Mig        string `short:"m" help:"Migrate script to a Go module directory."`
	Parallel   int    `short:"j" help:"Run pipe loop concurrently with N goroutines." default:"0"`
	Explain    bool   `help:"Print generated Go source without compiling."`
	NoCache    bool   `help:"Skip cache lookup; always recompile."`
	ClearCache bool   `help:"Delete the goscript cache directory and exit."`

	Script string   `arg:"" optional:"" help:"Inline Go code snippet or path to a Go script."`
	Args   []string `arg:"" optional:"" help:"Args passed to the script."`
}

var cli cliArgs

func ParseArgs(argv []string) Config {
	ctx := kong.Parse(&cli, kong.Name("goscript"))

	ok, errst := validate(cli)
	if !ok {
		ctx.Errorf("%s", errst)
		ctx.PrintUsage(false)
		os.Exit(1)
	}

	cfg := Config{
		Args:       cli.Args,
		OutputPath: cli.Out,
		MigrateDir: cli.Mig,
		FieldSep:   cli.FieldSep,
		Parallel:   cli.Parallel,
		NoCache:    cli.NoCache,
	}

	if isFilePath(cli.Script) {
		cfg.Input = InputScript
		cfg.ScriptPath = cli.Script
	} else {
		cfg.Input = InputInline
		cfg.InlineCode = cli.Script
	}

	if cli.FieldSep != "" && cfg.Input != InputInline {
		ctx.Errorf("-F only applies to inline mode")
		os.Exit(1)
	}
	if cli.Parallel > 0 && cfg.Input != InputInline {
		ctx.Errorf("-j only applies to inline mode")
		os.Exit(1)
	}

	switch {
	case cli.ClearCache:
		cfg.Action = ActionClearCache
	case cli.Explain:
		cfg.Action = ActionExplain
	case cli.Mig != "":
		cfg.Action = ActionMigrate
	case cli.Out != "":
		cfg.Action = ActionBuild
	default:
		cfg.Action = ActionRun
	}

	return cfg
}

func validate(cli cliArgs) (bool, string) {
	if cli.ClearCache {
		return true, ""
	}
	if cli.Script == "" {
		return false, "provide an inline code snippet or a script path"
	}
	if cli.Out != "" && cli.Mig != "" {
		return false, "can't have both output bin and migrate targets"
	}
	return true, ""
}

// isFilePath reports whether s refers to an existing file on disk.
// Used to distinguish inline code snippets from script paths.
func isFilePath(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}
