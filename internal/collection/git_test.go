package collection

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectGitInstalled(t *testing.T) {
	assert.True(t, DetectGitInstalled())
}

func TestLoadGitRepoInfo(t *testing.T) {
	info, err := LoadGitRepoInfo("../..") // this repository
	assert.NoError(t, err)
	assert.NotEmpty(t, info.Path)
	assert.Regexp(t, `^[0-9a-f]+\s*$`, info.HeadCommit)
	assert.NotEmpty(t, info.Status)
}

func TestExecGitCommand(t *testing.T) {
	output, err := ExecGitCommand(os.TempDir(), "--version")

	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
