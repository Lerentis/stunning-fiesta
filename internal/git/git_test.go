package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo, err := NewRepository("testrepo")
	if err != nil {
		t.Fatalf("NewRepository failed: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(repo.Path))

	// Check if .git directory exists
	gitDir := filepath.Join(repo.Path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf(".git directory does not exist in new repository")
	}
}

func TestChangeBranch(t *testing.T) {
	repo, err := NewRepository("branchrepo")
	if err != nil {
		t.Fatalf("NewRepository failed: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(repo.Path))

	if err := repo.ChangeBranch("feature/test"); err != nil {
		t.Errorf("ChangeBranch failed: %v", err)
	}
}

func TestCommitChanges(t *testing.T) {
	repo, err := NewRepository("commitrepo")
	if err != nil {
		t.Fatalf("NewRepository failed: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(repo.Path))

	// Create a dummy file
	dummyFile := filepath.Join(repo.Path, "dummy.txt")
	if err := os.WriteFile(dummyFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	if err := repo.CommitChanges("initial commit"); err != nil {
		t.Errorf("CommitChanges failed: %v", err)
	}
}

func TestCloneRepository(t *testing.T) {
	repo, err := CloneRepository("https://github.com/Lerentis/stunning-fiesta.git")
	if err != nil {
		t.Skipf("CloneRepository skipped (network or repo issue): %v", err)
		return
	}
	defer os.RemoveAll(filepath.Dir(repo.Path))

	// Check if .git directory exists
	gitDir := filepath.Join(repo.Path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf(".git directory does not exist in cloned repository")
	}
}
