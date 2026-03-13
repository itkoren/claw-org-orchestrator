package brain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func WriteTaskArtifact(repoDir, taskID, filename, content string) error {
	dir := filepath.Join(repoDir, "workspaces", "tasks", taskID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o644)
}

func Sync(repoDir string) error {
	// v1: best-effort
	steps := [][]string{
		{"git", "-C", repoDir, "pull", "--rebase"},
		{"git", "-C", repoDir, "add", "-A"},
		{"git", "-C", repoDir, "commit", "-m", "brain sync update", "--allow-empty"},
		{"git", "-C", repoDir, "push"},
	}
	for _, s := range steps {
		cmd := exec.Command(s[0], s[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%v failed: %w", s, err)
		}
	}
	return nil
}
