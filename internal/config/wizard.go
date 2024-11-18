package config

import (
	"bufio"
	"fmt"
	"github.com/NETWAYS/support-collector/modules/icinga2/icingaapi"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

const interactiveHelpText = `Welcome to the support-collector wizard!
We will guide you through all required details.

If you do not want to use the wizard, you can also pass an answer file containing the configuration.
For more details have a look at the official repository.
https://github.com/NETWAYS/support-collector`

var (
	validBooleanInputs = map[string]bool{
		"y":     true,
		"yes":   true,
		"true":  true,
		"n":     false,
		"no":    false,
		"false": false,
	}
)

type argument struct {
	name          string
	inputFunction func()
	dependency    func() bool
}

type Wizard struct {
	Scanner   *bufio.Scanner
	Arguments []argument
}

// NewWizard creates a new Wizard
func NewWizard() Wizard {
	return Wizard{
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

// Parse starts the interactive wizard for all Arguments defined in Wizard
func (w *Wizard) Parse(availableModules string) {
	// Print "welcome" text for wizard
	fmt.Printf("%s\n\nThe following modules are available:\n%s\n\n", interactiveHelpText, availableModules)

	for _, arg := range w.Arguments {
		if arg.dependency == nil {
			arg.inputFunction()
			continue
		}

		if ok := arg.dependency(); ok {
			arg.inputFunction()
			continue
		}
	}
}

// AddStringVar adds argument for a string variable
//
//	callback: Variable to save the input to
//	name: Internal name
//	defaultValue: Default
//	usage: usage string
//	required: bool
//	dependency: Add dependency function to validate if that argument will be added or not
func (w *Wizard) AddStringVar(callback *string, name, defaultValue, usage string, required bool, dependency func() bool) {
	arg := argument{
		name:       name,
		dependency: dependency,
	}

	arg.inputFunction = func() {
		if *callback != "" {
			defaultValue = *callback
		}

		w.newStringPromptWithDefault(callback, defaultValue, usage, required)
	}

	w.Arguments = append(w.Arguments, arg)
}

// AddSliceVarFromString reads a single string from stdin. This string will be separated by ',' and the resulting slice will be returned
//
//	callback: Variable to save the input to
//	name: Internal name
//	defaultValue: Default
//	usage: usage string
//	required: bool
//	dependency: Add dependency function to validate if that argument will be added or not
func (w *Wizard) AddSliceVarFromString(callback *[]string, name string, defaultValue []string, usage string, required bool, dependency func() bool) {
	arg := argument{
		name:       name,
		dependency: dependency,
	}

	arg.inputFunction = func() {
		if len(*callback) > 0 {
			defaultValue = *callback
		}

		var input string

		w.newStringPromptWithDefault(&input, strings.Join(defaultValue, ","), usage, required)
		*callback = strings.Split(strings.ReplaceAll(input, " ", ""), ",")
	}

	w.Arguments = append(w.Arguments, arg)
}

// AddBoolVar adds argument for a boolean variable
//
//	callback: Variable to save the input to
//	name: Internal name
//	defaultValue: Default
//	usage: usage string
//	dependency: Add dependency function to validate if that argument will be added or not
func (w *Wizard) AddBoolVar(callback *bool, name string, defaultValue bool, usage string, dependency func() bool) {
	arg := argument{
		name:       name,
		dependency: dependency,
	}

	arg.inputFunction = func() {
		w.newBoolPrompt(callback, defaultValue, usage)
	}

	w.Arguments = append(w.Arguments, arg)
}

// AddStringSliceVar adds argument for a slice of strings.
//
//	callback: Variable to save the input to
//	name: Internal name
//	defaultValue: Should that slice be collected via default?
//	initialPrompt: The initial stdout message to ask if you want to collect
//	inputPrompt: The recurring stdout message for each slice
//	dependency: Add dependency function to validate if that argument will be added or not
func (w *Wizard) AddStringSliceVar(callback *[]string, name string, defaultValue bool, initialPrompt string, inputPrompt string, dependency func() bool) {
	arg := argument{
		name:       name,
		dependency: dependency,
	}

	arg.inputFunction = func() {
		var (
			collect bool
			inputs  []string
		)

		w.newBoolPrompt(&collect, defaultValue, initialPrompt)

		if !collect {
			return
		}

		for collect {
			var input string

			w.newStringPrompt(&input, inputPrompt, true)
			inputs = append(inputs, input)

			w.newBoolPrompt(&collect, false, "Collect more?")
		}

		*callback = inputs
	}

	w.Arguments = append(w.Arguments, arg)
}

func (w *Wizard) AddIcingaEndpoints(callback *[]icingaapi.Endpoint, name, usage string, dependency func() bool) {
	arg := argument{
		name:       name,
		dependency: dependency,
	}

	arg.inputFunction = func() {
		// Ask if endpoint should be added. If not, return
		var collect bool

		w.newBoolPrompt(&collect, true, usage)

		if !collect {
			return
		}

		var endpoints []icingaapi.Endpoint

		for collect {
			e := w.newIcinga2EndpointPrompt()
			endpoints = append(endpoints, e)

			w.newBoolPrompt(&collect, false, "Collect more Icinga 2 API endpoints?")
		}

		*callback = endpoints
	}

	w.Arguments = append(w.Arguments, arg)
}

// newStringPromptWithDefault creates a new stdout / stdin prompt for a string
func (w *Wizard) newStringPrompt(callback *string, usage string, required bool) {
	for {
		fmt.Printf("%s: ", usage)

		if w.Scanner.Scan() {
			input := w.Scanner.Text()

			switch {
			case input != "":
				*callback = input
				return
			case input == "" && !required:
				return
			}
		} else {
			if err := w.Scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			}
		}
	}
}

// newStringPromptWithDefault creates a new stdout / stdin prompt for a string with default a value
func (w *Wizard) newStringPromptWithDefault(callback *string, defaultValue, usage string, required bool) {
	for {
		fmt.Printf("%s - (Default: '%s'): ", usage, defaultValue)

		if w.Scanner.Scan() {
			input := w.Scanner.Text()

			switch {
			case input != "":
				*callback = input
				return
			case input == "" && defaultValue != "":
				*callback = defaultValue
				return
			case input == "" && !required:
				return
			}
		} else {
			if err := w.Scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			}
		}
	}
}

// newIntPromptWithDefault creates a new stdout / stdin prompt for an int
func (w *Wizard) newIntPromptWithDefault(callback *int, defaultValue int, usage string, required bool) {
	for {
		fmt.Printf("%s - (Default: '%d'): ", usage, defaultValue)

		if w.Scanner.Scan() {
			input := w.Scanner.Text()

			switch {
			case input != "":
				converted, err := strconv.Atoi(input)
				if err != nil {
					logrus.Fatalf("could not convert '%s' to integer: %s", input, err)
				}

				*callback = converted

				return
			case input == "" && required:
				*callback = defaultValue
				return
			}
		} else {
			if err := w.Scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			}
		}
	}
}

// newBoolPrompt creates a new stdout / stdin prompt for a boolean
func (w *Wizard) newBoolPrompt(callback *bool, defaultValue bool, usage string) {
	for {
		fmt.Printf("%s [y/n] - (Default: '%t'): ", usage, defaultValue)

		if w.Scanner.Scan() {
			input := strings.ToLower(w.Scanner.Text())

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

// isValidBoolString validate the inputs for boolean. Valid are everything from var validBoolString
func isValidBoolString(input string) bool {
	if _, ok := validBooleanInputs[input]; !ok {
		return false
	}

	return true
}

// newIcinga2EndpointPrompt creates all needed stdout / stdin prompts to build an Icinga API endpoint. Returns icingaapi.Endpoint
func (w *Wizard) newIcinga2EndpointPrompt() (e icingaapi.Endpoint) {
	w.newStringPromptWithDefault(&e.Address, "127.0.0.1", "Host address / FQDN of the endpoint", true)
	w.newIntPromptWithDefault(&e.Port, 5665, "Port number of the endpoint", true) //nolint:mnd
	w.newStringPrompt(&e.Username, "Username for the api connection", true)
	w.newStringPrompt(&e.Password, "Password for the api connection", true)

	return e
}
