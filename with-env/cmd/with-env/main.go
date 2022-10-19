package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ingcr3at1on/x/sigctx"
	env "github.com/ingcr3at1on/x/with-env"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	rootFK     = `root`
	overloadFK = `overload`
	printFK    = `print`
)

var (
	rootF     *string
	overloadF *bool
	printF    *bool

	root = &cobra.Command{
		Use:           "with-env env-alias <cmd args>...",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if *printF {
				err = cobra.MinimumNArgs(1)(cmd, args)
			} else {
				err = cobra.MinimumNArgs(2)(cmd, args)
			}
			if err != nil {
				return err
			}

			return withEnv(cmd, args)
		},
	}
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	flags := root.PersistentFlags()
	rootF = flags.StringP(rootFK, "r", filepath.Join(home, ".env"), "root dir for env files")
	overloadF = flags.BoolP(overloadFK, "o", false, "overload existing environment values")
	printF = flags.BoolP(printFK, "p", false, "print env file instead of setting it's values")
}

func main() {
	if err := sigctx.StartWith(root.ExecuteContext); err != nil {
		log.Fatal(err)
	}
}

func withEnv(cmd *cobra.Command, args []string) error {
	var options []env.Option
	if *overloadF {
		options = append(options, env.OverloadOption)
	}
	if *printF {
		options = append(options, env.NewPrintOnlyOption(cmd.OutOrStdout()))
	}

	if err := env.LoadFrom(afero.NewOsFs(), *rootF, args[0], options...); err != nil {
		return err
	}
	if *printF {
		return nil
	}

	command := exec.CommandContext(cmd.Context(), args[1], args[2:]...)
	command.Stdin = cmd.InOrStdin()
	command.Stdout = cmd.OutOrStdout()
	command.Stderr = cmd.ErrOrStderr()
	return command.Run()
}
