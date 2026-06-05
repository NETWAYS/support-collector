package collection

import (
	"bytes"
	"os"
	"regexp"
	"testing"
)

func TestDetectGitInstalled(t *testing.T) {
	if !DetectGitInstalled() {
		t.Fatalf("expected git to be installed")
	}
}

func TestLoadGitRepoInfo(t *testing.T) {
	info, err := LoadGitRepoInfo("../..") // this repository

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if info.Path == "" {
		t.Error("expected Path not to be empty")
	}

	matched, err := regexp.MatchString(`^[0-9a-f]+\s*$`, info.HeadCommit)

	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}

	if !matched {
		t.Errorf("expected HeadCommit to match regex ^[0-9a-f]+\\s*$, got: %s", info.HeadCommit)
	}

	if info.Status == "" {
		t.Error("expected Status not to be empty")
	}
}

func TestExecGitCommand(t *testing.T) {
	output, err := ExecGitCommand(os.TempDir(), "--version")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if bytes.Equal(output, []byte{}) {
		t.Errorf("expected output not to be empty")
	}
}
