package env

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ingcr3at1on/x/with-env/internal"
	"github.com/joho/godotenv"
	"github.com/spf13/afero"
)

const (
	overloadOptionType  = "overload"
	printOnlyOptionType = "print-only"
)

type (
	Option interface {
		Type() string
	}

	option struct {
		optionType string
	}

	printOnlyOption struct {
		option
		out io.Writer
	}
)

var OverloadOption = option{
	optionType: overloadOptionType,
}

func (o option) Type() string {
	return o.optionType
}

// LoadFrom is a somewhat lazy wrapper around one or more hcl index files and godotenv.
// It leverages godotenv.Parse and simulates behavior to godotenv.Load/godotenv.Overload
func LoadFrom(afs afero.Fs, abs, alias string, options ...Option) error {
	r, err := internal.GetEnvFile(afs, abs, alias)
	if err != nil {
		return err
	}

	if r == nil {
		// TODO: loop through all env files for the complete list of aliases.
		return fmt.Errorf("environment alias %s not recognized", alias)
	}

	if opt, ok := hasOption(printOnlyOptionType, options...); ok {
		return internal.RedactWrite(opt.(printOnlyOption).out, r)
	}

	m, err := godotenv.Parse(r)
	if err != nil {
		return err
	}

	// Note this is a bit of a rip off of godotenv that I'd rather
	// not have, see about either opening a PR to get this specific
	// functionality in godotenv exposed in the API _or_ fork it
	// and use it as a basis for this (I'd rather not do that)...
	current := make(map[string]bool)
	raw := os.Environ()
	for _, ev := range raw {
		fields := strings.Split(ev, "=")
		current[fields[0]] = true
	}

	for k, v := range m {
		_, overload := hasOption(overloadOptionType, options...)
		if !current[k] || overload {
			os.Setenv(k, v)
		}
	}

	return nil
}

func NewPrintOnlyOption(out io.Writer) Option {
	return printOnlyOption{
		option: option{
			optionType: printOnlyOptionType,
		},
		out: out,
	}
}

func hasOption(optionType string, options ...Option) (Option, bool) {
	for _, opt := range options {
		if opt.Type() == optionType {
			return opt, true
		}
	}
	return option{}, false
}
