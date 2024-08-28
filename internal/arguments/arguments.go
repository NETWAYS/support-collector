package arguments

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

var (
	NonInteractive     bool
	validBooleanInputs = map[string]bool{
		"y":     true,
		"yes":   true,
		"n":     false,
		"no":    false,
		"true":  true,
		"false": false,
	}
)

const interactiveHelpText = `Welcome to the support-collector argument wizard!
We will guide you through all required details.

Available modules are: %s`

type Argument struct {
	Name          string
	InputFunction func()
	Dependency    func() bool
}

type Handler struct {
	scanner   *bufio.Scanner
	arguments []Argument
}

// New creates a new Handler object
func New() Handler {
	return Handler{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (args *Handler) CollectArgsFromStdin(availableModules string) {
	fmt.Printf(interactiveHelpText+"\n\n", availableModules)

	var errors []error

	for _, argument := range args.arguments {
		if argument.Dependency == nil {
			argument.InputFunction()
			continue
		}

		if ok := argument.Dependency(); ok {
			argument.InputFunction()
			continue
		}

		errors = append(errors, fmt.Errorf("%s is not matching the needed depenency", argument.Name))
	}
}

func (args *Handler) NewPromptStringVar(callback *string, name, defaultValue, usage string, required bool, dependency func() bool) {
	flag.StringVar(callback, name, defaultValue, usage)

	args.arguments = append(args.arguments, Argument{
		Name: name,
		InputFunction: func() {
			if *callback != "" {
				defaultValue = *callback
			}

			args.newStringPrompt(callback, defaultValue, usage, required)
		},
		Dependency: dependency,
	})
}

func (args *Handler) NewPromptStringSliceVar(callback *[]string, name string, defaultValue []string, usage string, required bool, dependency func() bool) {
	flag.StringSliceVar(callback, name, defaultValue, usage)

	args.arguments = append(args.arguments, Argument{
		Name: name,
		InputFunction: func() {
			if len(*callback) > 0 {
				defaultValue = *callback
			}

			var input string

			args.newStringPrompt(&input, strings.Join(defaultValue, ","), usage, required)
			*callback = strings.Split(input, ",")
		},
		Dependency: dependency,
	})
}

func (args *Handler) NewPromptBoolVar(callback *bool, name string, defaultValue bool, usage string, dependency func() bool) {
	flag.BoolVar(callback, name, defaultValue, usage)

	args.arguments = append(args.arguments, Argument{
		Name: name,
		InputFunction: func() {
			args.newBoolPrompt(callback, defaultValue, usage)
		},
		Dependency: dependency,
	})
}

func (args *Handler) newStringPrompt(callback *string, defaultValue, usage string, required bool) {
	for {
		fmt.Printf("%s - (Preselection: '%s'): ", usage, defaultValue)
		if args.scanner.Scan() {
			input := args.scanner.Text()
			if input != "" {
				*callback = input
				break
			} else if input == "" && defaultValue != "" {
				*callback = defaultValue
				break
			} else if input == "" && !required {
				break
			}
		} else {
			if err := args.scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
				break
			}
		}
	}

	return
}

func (args *Handler) newBoolPrompt(callback *bool, defaultValue bool, usage string) {
	for {
		fmt.Printf("%s [y/n] - (Preselection: '%t'): ", usage, defaultValue)

		if args.scanner.Scan() {
			input := strings.ToLower(args.scanner.Text())

			if input != "" && isValidBoolString(input) {
				*callback = validBooleanInputs[input]
				break
			} else if input == "" {
				*callback = defaultValue
				break
			}
		}
	}
}

func isValidBoolString(input string) bool {
	if _, ok := validBooleanInputs[input]; !ok {
		return false
	}

	return true
}
