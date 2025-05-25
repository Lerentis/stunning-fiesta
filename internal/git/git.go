package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Repository struct {
	Path string
}

func NewRepository(name string) (*Repository, error) {
	dir, err := os.MkdirTemp("", "stunning-fiesta-git-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	repoPath := filepath.Join(dir, name)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return nil, err
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git init failed: %v, output: %s", err, string(output))
	}
	return &Repository{Path: repoPath}, nil
}

func CloneRepository(url string) (*Repository, error) {
	dir, err := os.MkdirTemp("", "stunning-fiesta-git-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	cmd := exec.Command("git", "clone", url, dir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git clone failed: %v, output: %s", err, string(output))
	}
	return &Repository{Path: dir}, nil
}

func (r *Repository) ChangeBranch(name string) error {
	cmd := exec.Command("git", "checkout", "-B", name)
	cmd.Dir = r.Path
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (r *Repository) AddChanges() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = r.Path
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (r *Repository) CommitChanges(message string) error {
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = r.Path
	if output, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %v, output: %s", err, string(output))
	}
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = r.Path
	if output, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (r *Repository) PushChanges() error {
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = r.Path
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (r *Repository) PullChanges() error {
	cmd := exec.Command("git", "pull", "--rebase")
	cmd.Dir = r.Path
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, string(output))
	}
	return nil
}
