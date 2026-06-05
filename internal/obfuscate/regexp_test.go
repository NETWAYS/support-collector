package obfuscate

import (
	"testing"
)

func TestNewRegexpKeyValue(t *testing.T) {
	re := NewRegexpKeyValue("password")

	tests := []struct {
		input    string
		expected bool
	}{
		{`password = "test"`, true},
		{`  password = "test"`, true},
		{`;password = "test"`, true},
		{`  //password = "test"`, true},
		{`Password = "test"`, true},
		{`db_password = "test"`, true},
		{`user = "test"`, false},
		{`password`, false},
	}

	for _, test := range tests {
		matched := re.MatchString(test.input)
		if matched != test.expected {
			t.Errorf("input: %q, expected: %v, got: %v", test.input, test.expected, matched)
		}
	}
}

func TestNewExtensionMatch(t *testing.T) {
	re := NewExtensionMatch("ini")

	tests := []struct {
		input    string
		expected bool
	}{
		{`test.ini`, true},
		{`/etc/bla/test.ini`, true},
		{`test.txt`, false},
	}

	for _, test := range tests {
		matched := re.MatchString(test.input)
		if matched != test.expected {
			t.Errorf("input: %q, expected: %v, got: %v", test.input, test.expected, matched)
		}
	}
}

func TestNewCommandMatch(t *testing.T) {
	if !NewCommandMatch("icinga2", "daemon").MatchString("icinga2 daemon") {
		t.Errorf("Expected match on icinga2 daemon")
	}

	if !NewCommandMatch("icinga2").MatchString("icinga2") {
		t.Errorf("Expected match on icinga2")
	}
}
