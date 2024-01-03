package argument

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"slices"
	"strings"
)

// Argument is a predefined blueprint. With that you can define all needed metadata for the argument.
//
// The input will be returned to the variable given in Variable.
//
// Example:
//
//	argument.InitArgs([]argument.Argument{
//		{
//			Callback:    argument.NewStringPrompt,
//			Required:    false,
//			Description: "All modules are activated by default. If you would like to enable only explicit modules, please enter them as a comma-separated list. Example: 'icinga2,icingadb,webservers,base'",
//			InputPrompt: "Enable only single modules (default: all enabled)",
//			Variable:    &enabledModules,
//		},
//	})
type Argument struct {
	Callback    func(argument *Argument)
	Required    bool
	Description string
	Longhand    string
	InputPrompt string
	Variable    *interface{}
}

var (
	// boolOptions defines allowed inputs for true and false
	boolOptions []string

	// helpOptions defines allowed inputs to print the help
	helpOptions []string

	EnabledModules interface{}
)

// TODO
func InitArgs(args []Argument) {
	// TODO print small help introduction
	helpOptions = []string{"?", "h", "H"}
	boolOptions = []string{"Y", "N", "y", "n"}

	for _, arg := range args {
		arg.Callback(&arg)
	}
}

func returnArgs() []Argument {
	return []Argument{
		{
			Callback:    NewListPrompt,
			Required:    false,
			Description: "All modules are activated by default. If you would like to enable only explicit modules, please enter them as a comma-separated list. Example: 'icinga2,icingadb,webservers,base'",
			InputPrompt: "Enable only single modules (default: all enabled)",
			Variable:    &EnabledModules,
		},
	}
}

// TODO
func NewStringPrompt(arg *Argument) {
	var input string
	r := bufio.NewReader(os.Stdin)

	for {
		printBoldStdOut(arg.InputPrompt + ": ")

		input, _ = r.ReadString('\n')

		if input != "" && slices.Contains(helpOptions, strings.TrimSpace(input)) {
			fmt.Fprintln(os.Stdout, arg.Description)
			input = ""
			continue
		}

		if strings.TrimSpace(input) != "" || !arg.Required {
			break
		}
	}

	fmt.Println("res: ", input)
	*arg.Variable = strings.TrimSpace(input)
}

// TODO
func NewListPrompt(arg *Argument) {
}

// TODO
func NewDigitPrompt() {

}

// TODO
func NewBooleanPrompt(arg *Argument) {
	var s string
	possibleInputs := []string{"Y", "N", "y", "n"}

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stderr, arg.InputPrompt+": ")
		s, _ = r.ReadString('\n')
		if slices.Contains(possibleInputs, strings.TrimSpace(s)) {
			break
		}
	}

	*arg.Variable = strings.TrimSpace(s)
}

// printBoldStdOut takes string and prints text bold to stdout
//
//	Example: printBoldStdOut("This string is bold and will be printed to stdout")
func printBoldStdOut(text string) {
	out := color.New(color.Bold).FprintfFunc()
	out(os.Stdout, text)
}
