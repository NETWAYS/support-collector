package obfuscate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

const iniExample = `[default]
user = "test"
password = "very-secret"
`

const iniResult = `[default]
user = "test"
password = <HIDDEN>
`

func ExampleObfuscator() {
	o := New(KindFile, regexp.MustCompile(`\.ini$`), regexp.MustCompile(`password\s*=\s*(.*)`))

	content := []byte(`password = "secret"`)

	if o.IsAccepting(KindFile, "test.ini") {
		count, data, err := o.Process(content)
		fmt.Println(err)
		fmt.Println(count)
		fmt.Println(string(data))
	}

	// Output: <nil>
	// 1
	// password = <HIDDEN>
}

func TestReplacePattern(t *testing.T) {
	replacement, count := ReplacePattern(`password = "XXX"`, regexp.MustCompile(`password\s*=\s*(.*)`))
	assert.Equal(t, uint(1), count)
	assert.Equal(t, `password = <HIDDEN>`, replacement)

	replacement, count = ReplacePattern(`password = "XXX"`, regexp.MustCompile(`password\s*=.*`))
	assert.Equal(t, uint(1), count)
	assert.Equal(t, `<HIDDEN>`, replacement)
}

func TestObfuscator_IsAccepting(t *testing.T) {
	o := New(KindFile, NewExtensionMatch("ini"), regexp.MustCompile(`^$`))

	assert.True(t, o.IsAccepting(KindFile, "test.ini"))
	assert.False(t, o.IsAccepting(KindFile, "test.txt"))
	assert.False(t, o.IsAccepting(KindOutput, "echo"))

	o.Kind = KindAny
	assert.True(t, o.IsAccepting(KindOutput, "test.ini"))

	o.Kind = KindOutput
	o.WithAffecting(regexp.MustCompile(`^icinga2 daemon -C`))

	assert.True(t, o.IsAccepting(KindOutput, "icinga2 daemon -C"))
	assert.False(t, o.IsAccepting(KindOutput, "icinga2 daemon"))
}

func TestObfuscator_Process(t *testing.T) {
	o := &Obfuscator{
		Replacements: []*regexp.Regexp{
			regexp.MustCompile(`^password\s*=\s*(.*)`),
		},
	}

	count, out, err := o.Process([]byte("default content\r\n"))
	assert.NoError(t, err)
	assert.Equal(t, uint(0), count)
	assert.Equal(t, "default content\r\n", string(out))

	count, out, err = o.Process([]byte(iniExample))
	assert.NoError(t, err)
	assert.Equal(t, uint(1), count)
	assert.Equal(t, iniResult, string(out))
}

func TestNewFile(t *testing.T) {
	o := NewFile(`^password\s*=\s*(.*)`, "conf")
	assert.Equal(t, KindFile, o.Kind)
	assert.Len(t, o.Affecting, 1)
	assert.Len(t, o.Replacements, 1)
}

func TestNewOutput(t *testing.T) {
	o := NewOutput(`^password\s*=\s*(.*)`, "icinga2", "daemon", "-C")
	assert.Equal(t, KindOutput, o.Kind)
	assert.Len(t, o.Affecting, 1)
	assert.Len(t, o.Replacements, 1)

	assert.True(t, o.IsAccepting(KindOutput, "icinga2 daemon -C"))
	assert.False(t, o.IsAccepting(KindOutput, "icinga2 daemon"))
}
