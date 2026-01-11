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
)

type Config struct {
	Input  InputMode
	Action Action

	ScriptPath string
	Args       []string
	InlineCode string

	OutputPath string
	MigrateDir  string
}

type cliArgs struct {
	Code string `short:"c" help:"Inline Go code."`
	Out  string `short:"o" help:"Build output path. Not Implemented."`
	Mig  string `short:"m" help:"Migrate target dir. Not Implemented."`

	Script string   `arg:"" optional:"" help:"Go script path (default run)."`
	Args   []string `arg:"" optional:"" help:"Args for script."`
}

var cli cliArgs

func ParseArgs(argv []string) Config {
	ctx := kong.Parse(&cli, kong.Name("goscript"))

	ok, errst := validate(cli)
	if (!ok) {
		ctx.Errorf("%s", errst)
		ctx.PrintUsage(false);
		os.Exit(1)
	}

	cfg := Config{
		ScriptPath: cli.Script,
		Args:       cli.Args,
		InlineCode: cli.Code,
		OutputPath: cli.Out,
		MigrateDir:  cli.Mig,
	}

	if cli.Code != "" {
		cfg.Input = InputInline
	} else {
		cfg.Input = InputScript
	}

	if cli.Mig != "" {
		cfg.Action = ActionMigrate
	} else if cli.Out != "" {
		cfg.Action = ActionBuild
	} else {
		cfg.Action = ActionRun
	}

	return cfg
}

func validate(cli cliArgs) (bool, string) {
	if ((cli.Code != "") == (cli.Script != "")) {
		return false, "Only one of inline code or input script should be set";
	}

	if (cli.Out != "" && cli.Mig != "") {
		return false, "Can't have both output bin and migrate targets";
	}

	return true, ""
}
