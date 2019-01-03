package unit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// derived from https://github.com/docker/docker/blob/master/pkg/mflag/flag.go
var flagNoValue = map[string]bool{
	"-rm": true, "rm": true, "t": true, "i": true,
}

// ErrHelp is the error returned if the flag -help is invoked but no such flag is defined.
var ErrHelp = errors.New("flag: help requested")

// ErrRetry is the error returned if you need to try letter by letter
var ErrRetry = errors.New("flag: retry")

type FlagSet struct {
	parsed bool
	actual map[string]string
	args   []string          // arguments after flags
	output io.Writer         // nil means stderr; use Out() accessor
	env    map[string]string // enviornment variables specified with -e
}

func trimQuotes(str string) string {
	if len(str) == 0 {
		return str
	}
	type quote struct {
		start, end byte
	}

	// All valid quote types.
	quotes := []quote{
		// Double quotes
		{
			start: '"',
			end:   '"',
		},

		// Single quotes
		{
			start: '\'',
			end:   '\'',
		},
	}

	for _, quote := range quotes {
		// Only strip if outermost match.
		if str[0] == quote.start && str[len(str)-1] == quote.end {
			str = str[1 : len(str)-1]
			break
		}
	}

	return str
}

// Out returns the destination for usage and error messages.
func (f *FlagSet) Out() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

// parseOne parses one flag. It reports whether a flag was seen.
func (f *FlagSet) parseOne() (bool, string, error) {
	if len(f.args) == 0 {
		return false, "", nil
	}
	s := f.args[0]
	if len(s) == 0 || s[0] != '-' || len(s) == 1 {
		return false, "", nil
	}
	if s[1] == '-' && len(s) == 2 { // "--" terminates the flags
		f.args = f.args[1:]
		return false, "", nil
	}
	name := s[1:]
	if len(name) == 0 || name[0] == '=' {
		return false, "", fmt.Errorf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	has_value := false
	value := ""
	if i := strings.Index(name, "="); i != -1 {
		value = trimQuotes(name[i+1:])
		has_value = true
		name = name[:i]
	}

	// It must have a value, which might be the next argument.
	if !has_value && len(f.args) > 0 && !flagNoValue[name] {
		// value is the next arg
		if len(f.args[0]) > 0 && f.args[0][0] != '-' {
			has_value = true
			value, f.args = f.args[0], f.args[1:]
		}
	}

	// Store environment variables separately
	if name == "e" || name == "-env" {
		if i := strings.Index(value, "="); i == -1 {
			f.env[value] = ""
		} else {
			f.env[value[:i]] = value[i+1:]
		}
	} else {
		f.actual[name] = value
	}

	return true, "", nil
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	f.actual = map[string]string{}
	f.env = map[string]string{}
	for {
		seen, name, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		if err == ErrRetry {
			if len(name) > 1 {
				err = nil
				for _, letter := range strings.Split(name, "") {
					f.args = append([]string{"-" + letter}, f.args...)
					seen2, _, err2 := f.parseOne()
					if seen2 {
						continue
					}
					if err2 != nil {
						// err = fmt.Errorf("flag provided but not defined: -%s", name)
						continue
						break
					}
				}
				if err == nil {
					continue
				}
			} else {
				// err = fmt.Errorf("flag provided but not defined: -%s", name)
				continue
			}
		}
		return err
	}
	return nil
}

func (f *FlagSet) Args() []string {
	return f.args
}

func (f *FlagSet) Values() map[string]string {
	return f.actual
}

func (f *FlagSet) Env() map[string]string {
	return f.env
}
