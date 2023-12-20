package collection

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GitRepoInfo struct {
	Path       string
	HeadCommit string
	Remotes    string
	Describe   string
	Status     string
}

func DetectGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func IsGitRepository(path string) (string, bool) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}

	// Check if git dir exists
	stat, err := os.Stat(filepath.Join(path, ".git"))
	if err != nil || !stat.IsDir() {
		return absPath, false
	}

	return absPath, true
}

func LoadGitRepoInfo(path string) (*GitRepoInfo, error) {
	absPath, ok := IsGitRepository(path)
	if !ok {
		return nil, fmt.Errorf("not a git repository: %s", absPath) //nolint:goerr113
	}

	i := &GitRepoInfo{
		Path: absPath,
	}

	// Collecting various state information via commands
	for v, cmd := range map[*string][]string{
		&i.HeadCommit: {"log", "-1", `--format=%H`},
		&i.Remotes:    {"remote", "-v"},
		&i.Describe:   {"describe", "--tags"},
		&i.Status:     {"status"},
	} {
		output, err := ExecGitCommand(absPath, cmd...)

		*v = string(output)
		if err != nil {
			*v += "\n" + err.Error()
		}
	}

	return i, nil
}

func ExecGitCommand(dir string, command ...string) ([]byte, error) {
	arguments := []string{"--git-dir", filepath.Join(dir, ".git"), "--work-tree", dir}
	arguments = append(arguments, command...)

	return LoadCommandOutput("git", arguments...)
}
