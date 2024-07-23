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
		"y":   true,
		"yes": true,
		"n":   false,
		"no":  false,
	}
)

type Handler struct {
	scanner *bufio.Scanner
	prompts []func()
}

func NewHandler() Handler {
	return Handler{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (h *Handler) ReadArgumentsFromStdin() {
	for _, prompt := range h.prompts {
		prompt()
	}
}

func (h *Handler) NewPromptStringVar(callback *string, name, defaultValue, usage string, required bool) {
	flag.StringVar(callback, name, defaultValue, usage)

	h.prompts = append(h.prompts, func() {
		if *callback != "" {
			defaultValue = *callback
		}

		h.newStringPrompt(callback, defaultValue, usage, required)
	})
}

func (h *Handler) NewPromptStringSliceVar(callback *[]string, name string, defaultValue []string, usage string, required bool) {
	flag.StringSliceVar(callback, name, defaultValue, usage)

	h.prompts = append(h.prompts, func() {
		if len(*callback) > 0 {
			defaultValue = *callback
		}

		var input string

		h.newStringPrompt(&input, strings.Join(defaultValue, ","), usage, required)
		*callback = strings.Split(input, ",")
	})
}

func (h *Handler) NewPromptBoolVar(callback *bool, name string, defaultValue bool, usage string) {
	flag.BoolVar(callback, name, defaultValue, usage)

	h.prompts = append(h.prompts, func() {
		h.newBoolPrompt(callback, defaultValue, usage)
	})
}

func (h *Handler) newStringPrompt(callback *string, defaultValue, usage string, required bool) {
	for {
		fmt.Printf("%s - (Preselection: '%s'): ", usage, defaultValue)
		if h.scanner.Scan() {
			input := h.scanner.Text()
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
			if err := h.scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
				break
			}
		}
	}

	return
}

func (h *Handler) newBoolPrompt(callback *bool, defaultValue bool, usage string) {
	for {
		fmt.Printf("%s [y/n] - (Preselection: '%t'): ", usage, defaultValue)

		if h.scanner.Scan() {
			input := strings.ToLower(h.scanner.Text())

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
