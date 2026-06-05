package obfuscate

import (
	"fmt"
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
		count, data, err := o.Process(content, "")
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

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	if replacement != `password = <HIDDEN>` {
		t.Errorf("expected replacement 'password = <HIDDEN>', got %q", replacement)
	}

	replacement, count = ReplacePattern(`password = "XXX"`, regexp.MustCompile(`password\s*=.*`))

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	if replacement != `<HIDDEN>` {
		t.Errorf("expected replacement '<HIDDEN>', got %q", replacement)
	}
}

func TestObfuscator_IsAccepting(t *testing.T) {
	o := New(KindFile, NewExtensionMatch("ini"), regexp.MustCompile(`^$`))

	if !o.IsAccepting(KindFile, "test.ini") {
		t.Fatalf("expected isAccepting to be true")
	}

	if o.IsAccepting(KindFile, "test.txt") {
		t.Fatalf("expected isAccepting to be false")
	}

	if o.IsAccepting(KindFile, "echo") {
		t.Fatalf("expected isAccepting to be false")
	}

	o.Kind = KindAny

	if !o.IsAccepting(KindFile, "test.ini") {
		t.Fatalf("expected isAccepting to be true")
	}

	o.Kind = KindOutput
	o.WithAffecting(regexp.MustCompile(`^icinga2 daemon -C`))

	if !o.IsAccepting(KindOutput, "icinga2 daemon -C") {
		t.Fatalf("expected isAccepting to be true")
	}

	if o.IsAccepting(KindOutput, "icinga2 daemon") {
		t.Fatalf("expected isAccepting to be false")
	}
}

func TestObfuscator_Process(t *testing.T) {
	o := &Obfuscator{
		ReplacePattern: regexp.MustCompile(`^password\s*=\s*(.*)`),
	}

	count, out, err := o.Process([]byte("default content\r\n"), "")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if count != uint(0) {
		t.Errorf("expected count to be 0")
	}

	if "default content\r\n" != string(out) {
		t.Errorf("expected %v, got %v", "default content", string(out))
	}

	count, out, err = o.Process([]byte(iniExample), "")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if count != uint(1) {
		t.Errorf("expected count to be 1")
	}

	if iniResult != string(out) {
		t.Errorf("expected %v, got %v", iniResult, string(out))
	}
}

func TestNewFile(t *testing.T) {
	o := NewFile(`^password\s*=\s*(.*)`, "conf")

	if o.Kind != KindFile {
		t.Errorf("expected Kind to be %v, got %v", KindFile, o.Kind)
	}

	if len(o.ShouldAffect) != 1 {
		t.Errorf("expected ShouldAffect length 1, got %d", len(o.ShouldAffect))
	}

	if len(o.ShouldAffect) == 0 {
		t.Error("expected ShouldAffect to be non-empty")
	}
}

func TestNewOutput(t *testing.T) {
	o := NewOutput(`^password\s*=\s*(.*)`, "icinga2", "daemon", "-C")

	if o.Kind != KindOutput {
		t.Errorf("expected Kind to be %v, got %v", KindOutput, o.Kind)
	}

	if len(o.ShouldAffect) != 1 {
		t.Errorf("expected ShouldAffect length 1, got %d", len(o.ShouldAffect))
	}

	if o.ReplacePattern == nil {
		t.Error("expected ReplacePattern to be non-empty")
	}

	if !o.IsAccepting(KindOutput, "icinga2 daemon -C") {
		t.Fatalf("expected isAccepting to be true")
	}

	if o.IsAccepting(KindOutput, "icinga2 daemon") {
		t.Fatalf("expected isAccepting to be true")
	}
}
