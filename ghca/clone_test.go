package ghca

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestNewCloner(t *testing.T) {
	c := NewCloner("/path/to/dest", nil, true, true)
	if c.git != "git" {
		t.Error("Git command should be initialized as 'git' by default:", c.git)
	}
	if c.dest != "/path/to/dest" {
		t.Error("Distination to clone should be set to given path:", c.dest)
	}
	if !c.deep {
		t.Error("deep clone should be set:", c.deep)
	}
	if !c.ssh {
		t.Error("ssh should be set:", c.ssh)
	}

	os.Setenv("GIT_EXECUTABLE_PATH", "/path/to/git")
	c = NewCloner("/path/to/dest", nil, false, false)
	if c.git != "/path/to/git" {
		t.Error("Git command should respect environment variable $GIT_EXECUTABLE_PATH:", c.git)
	}

	os.Setenv("GIT_EXECUTABLE_PATH", "")
}

func testRepos(repos []string, para int, t *testing.T) {
	c := NewCloner("test", nil, false, false)
	defer func() {
		os.RemoveAll("test")
	}()
	c.Err = make(chan error, 10)
	c.Start(para)

	go func() {
		for err := range c.Err {
			t.Error("Error reported from cloner:", err)
		}
	}()

	for _, r := range repos {
		c.Clone(r)
	}
	c.Shutdown()

	for _, r := range repos {
		p := filepath.FromSlash("test/" + r)
		s, err := os.Stat(p)
		if err != nil {
			t.Fatal("Cloned directory not found:", p, err)
		}
		if !s.IsDir() {
			t.Fatal("It should clone directory", p)
		}
	}
}

func TestClone1Repo(t *testing.T) {
	testRepos([]string{"rhysd/github-complete.vim"}, 0, t)
}

func TestCloneAFewRepos(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	repos := []string{
		"rhysd/clever-f.vim",
		"rhysd/neovim-component",
		"rhysd/vim-gfm-syntax",
	}
	testRepos(repos, 0, t)
}

func TestCloneManyRepos(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	repos := []string{
		"rhysd/inu-snippets",
		"rhysd/conflict-marker.vim",
		"rhysd/committia.vim",
		"rhysd/vim-dachs",
		"rhysd/rust-doc.vim",
		"vim-crystal/vim-crystal",
		"rhysd/vim-wasm",
		"rhysd/unite-go-import.vim",
		"rhysd/NyaoVim",
		"rhysd/vim-color-spring-night",
	}
	testRepos(repos, 0, t)
}

func TestCloneWithExtract(t *testing.T) {
	re := regexp.MustCompile("\\.vim$")
	c := NewCloner("test", re, false, false)
	defer func() {
		os.RemoveAll("test")
	}()
	c.Err = make(chan error, 10)
	c.Start(0)

	go func() {
		for err := range c.Err {
			t.Error("Error reported from cloner:", err)
		}
	}()

	c.Clone("rhysd/clever-f.vim")
	c.Shutdown()

	if err := filepath.Walk("test/rhysd/clever-f.vim", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && !re.MatchString(path) {
			t.Error("File not matching to 'extract' remains:", path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCloneNotExistingRepo(t *testing.T) {
	c := NewCloner("test", nil, false, false)
	c.Err = make(chan error, 10)
	c.Start(0)

	c.Clone("")
	c.Shutdown()

	select {
	case err, ok := <-c.Err:
		if !ok || err == nil {
			t.Fatal("Error not reported")
		}
	default:
		t.Fatal("Error not reported")
	}
}

func TestClone1Worker(t *testing.T) {
	testRepos([]string{"rhysd/github-complete.vim"}, 1, t)
}

func TestShallowClone(t *testing.T) {
	c := NewCloner("test", nil, true, false)
	defer func() {
		os.RemoveAll("test")
	}()
	c.Err = make(chan error, 10)
	c.Start(0)

	go func() {
		for err := range c.Err {
			t.Error("Error reported from cloner:", err)
		}
	}()

	c.Clone("rhysd/cargo-husky")
	c.Shutdown()

	cmd := exec.Command("git", "log", "--oneline")
	cmd.Dir = filepath.Join("test", "rhysd", "cargo-husky")
	bytes, err := cmd.Output()
	if err != nil {
		t.Fatal("git log failed:", err)
	}
	lines := strings.Split(string(bytes), "\n")
	lines = lines[:len(lines)-1]
	if len(lines) == 1 {
		t.Fatal("Log should not be one line since it's deep clone:", lines)
	}
}

func TestCloneSSH(t *testing.T) {
	c := NewCloner("test", nil, false, true)
	c.Err = make(chan error, 10)
	c.Start(0)

	c.Clone("rhysd/github-clone-all")
	c.Shutdown()

	url := "git@github.com:rhysd/github-clone-all.git"

	select {
	case err, ok := <-c.Err:
		if !ok || err != nil {
			t.Fatal("Error channel is broken:", ok, err)
		}
		msg := err.Error()
		if !strings.Contains(msg, url) {
			t.Error("Unexpected error:", msg)
		}
	default:
		cmd := exec.Command("git", "config", "--get", "remote.origin.url")
		cmd.Dir = filepath.Join("test", "rhysd", "github-clone-all")
		bytes, err := cmd.Output()
		if err != nil {
			t.Fatal("git config --get failed:", err)
		}
		out := string(bytes)

		if !strings.HasPrefix(out, url) {
			t.Error("Not a SSH URL", out)
		}
	}
}
