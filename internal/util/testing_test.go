package util

import (
	"testing"

	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

var (
	noObfuscator []*obfuscate.Obfuscator
)

func testObfuscator() []*obfuscate.Obfuscator {
	return []*obfuscate.Obfuscator{
		obfuscate.NewOutput(`.*`, "command"),
		obfuscate.NewFile(`.*`, "log"),
	}
}

func TestLoadTestdata(t *testing.T) {
	mockT := new(testing.T)

	actual := LoadTestdata(mockT, "command.txt")

	if actual != "output\n" {
		t.Error("\nActual: ", actual, "\nExpected: ", "output\n")
	}

	if mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be false")
	}

	LoadTestdata(mockT, "nonexisting.txt")

	if !mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be true")
	}
}

func TestAssertObfuscation(t *testing.T) {
	mockT := new(testing.T)
	AssertObfuscation(mockT, noObfuscator, obfuscate.KindFile, "a", "b", "c")

	if !mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be true")
	}

	o := testObfuscator()

	mockT = new(testing.T)
	AssertObfuscation(mockT, o, obfuscate.KindOutput, "command", "b", "<HIDDEN>")
	AssertObfuscation(mockT, o, obfuscate.KindFile, "test.log", "b", "<HIDDEN>")

	if mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be false")
	}
}

func TestAssertObfuscationExample(t *testing.T) {
	mockT := new(testing.T)
	o := testObfuscator()

	AssertObfuscationExample(mockT, o, obfuscate.KindOutput, "command")
	AssertObfuscationExample(t, o, obfuscate.KindFile, "file/test.log")

	if mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be false")
	}
}

func TestAssertAllObfuscatorsTested(t *testing.T) {
	mockT := new(testing.T)

	AssertAllObfuscatorsTested(mockT, noObfuscator)

	if !mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be true")
	}

	mockT = new(testing.T)
	o := testObfuscator()

	AssertObfuscation(mockT, o, obfuscate.KindOutput, "command", "b", "<HIDDEN>")
	AssertObfuscation(mockT, o, obfuscate.KindFile, "test.log", "b", "<HIDDEN>")
	AssertAllObfuscatorsTested(t, o)

	if mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be false")
	}

	mockT = new(testing.T)

	AssertAllObfuscatorsTested(mockT, testObfuscator())

	if !mockT.Failed() {
		t.Fatalf("expected mockT.Failed to be false")
	}
}
