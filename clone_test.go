package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestNewCloner(t *testing.T) {
	c := newCloner("/path/to/dest", nil)
	if c.git != "git" {
		t.Error("Git command should be initialized as 'git' by default:", c.git)
	}
	if c.dest != "/path/to/dest" {
		t.Error("Distination to clone should be set to given path:", c.dest)
	}

	os.Setenv("GIT_EXECUTABLE_PATH", "/path/to/git")
	c = newCloner("/path/to/dest", nil)
	if c.git != "/path/to/git" {
		t.Error("Git command should respect environment variable $GIT_EXECUTABLE_PATH:", c.git)
	}

	os.Setenv("GIT_EXECUTABLE_PATH", "")
}

func testRepos(repos []string, t *testing.T) {
	c := newCloner("test", nil)
	defer func() {
		os.RemoveAll("test")
	}()
	c.err = make(chan error, 10)
	c.start()

	go func() {
		for err := range c.err {
			t.Error("Error reported from cloner:", err)
		}
	}()

	for _, r := range repos {
		c.clone(r)
	}
	c.shutdown()

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
	testRepos([]string{"rhysd/github-complete.vim"}, t)
}

func TestCloneAFewRepos(t *testing.T) {
	repos := []string{
		"rhysd/clever-f.vim",
		"rhysd/neovim-component",
		"rhysd/vim-gfm-syntax",
	}
	testRepos(repos, t)
}

func TestCloneManyRepos(t *testing.T) {
	repos := []string{
		"rhysd/inu-snippets",
		"rhysd/conflict-marker.vim",
		"rhysd/committia.vim",
		"rhysd/vim-dachs",
		"rhysd/rust-doc.vim",
		"rhysd/vim-crystal",
		"rhysd/vim-wasm",
		"rhysd/unite-go-import.vim",
		"rhysd/NyaoVim",
		"rhysd/vim-color-spring-night",
	}
	testRepos(repos, t)
}

func TestCloneWithExtract(t *testing.T) {
	re := regexp.MustCompile("\\.vim$")
	c := newCloner("test", re)
	defer func() {
		os.RemoveAll("test")
	}()
	c.err = make(chan error, 10)
	c.start()

	go func() {
		for err := range c.err {
			t.Error("Error reported from cloner:", err)
		}
	}()

	c.clone("rhysd/clever-f.vim")
	c.shutdown()

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
	c := newCloner("test", nil)
	c.ssh = true
	c.err = make(chan error, 10)
	c.start()

	c.clone("rhysd/not-existing-repository")
	c.shutdown()

	select {
	case err, ok := <-c.err:
		if !ok || err == nil {
			t.Fatal("Error not reported")
		}
	default:
		t.Fatal("Error not reported")
	}
}
