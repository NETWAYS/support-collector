package arguments

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func NewStringVar(callback *string, value string, usage string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s - (Default: %s): ", usage, value)
		if scanner.Scan() {
			input := scanner.Text()
			if input != "" {
				*callback = input
				break
			} else {
				*callback = value
				break
			}
		} else {
			if err := scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
				break
			}
		}
	}
}

func NewStringSliceVar(callback *[]string, value []string, usage string) {
	var input string
	NewStringVar(&input, strings.Join(value, ","), usage)

	*callback = strings.Split(input, ",")
}

func NewBoolVar() {

}
