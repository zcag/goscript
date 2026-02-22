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
	Code       string `short:"c" help:"Inline Go code. Magic vars: x/line (current line), i/idx (index), f/fields (split fields), lines (all lines)."`
	FieldSep   string `short:"F" help:"Field separator for f/fields in pipe mode (default: whitespace split)."`
	Out        string `short:"o" help:"Build output path."`
	Mig        string `short:"m" help:"Migrate script to a Go module directory."`
	Parallel   int    `short:"j" help:"Run pipe loop concurrently with N goroutines." default:"0"`
	Explain    bool   `help:"Print generated Go source without compiling."`
	NoCache    bool   `help:"Skip cache lookup; always recompile."`
	ClearCache bool   `help:"Delete the goscript cache directory and exit."`

	Script string   `arg:"" optional:"" help:"Go script path (default run)."`
	Args   []string `arg:"" optional:"" help:"Args for script."`
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
		ScriptPath: cli.Script,
		Args:       cli.Args,
		InlineCode: cli.Code,
		OutputPath: cli.Out,
		MigrateDir: cli.Mig,
		FieldSep:   cli.FieldSep,
		Parallel:   cli.Parallel,
		NoCache:    cli.NoCache,
	}

	if cli.Code != "" {
		cfg.Input = InputInline
	} else {
		cfg.Input = InputScript
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
	if cli.Code == "" && cli.Script == "" {
		return false, "provide inline code with -c or a script path"
	}
	if cli.Code != "" && cli.Script != "" {
		return false, "only one of inline code or input script should be set"
	}
	if cli.Out != "" && cli.Mig != "" {
		return false, "can't have both output bin and migrate targets"
	}
	if cli.FieldSep != "" && cli.Code == "" {
		return false, "-F only applies to inline mode (-c)"
	}
	if cli.Parallel > 0 && cli.Code == "" {
		return false, "--parallel only applies to inline mode (-c)"
	}
	return true, ""
}
