package main

import (
	"os"
	"testing"
	"time"
)

func TestNewCloner(t *testing.T) {
	c := newCloner("/path/to/dist")
	if c.git != "git" {
		t.Error("Git command should be initialized as 'git' by default:", c.git)
	}
	if c.dist != "/path/to/dist" {
		t.Error("Distination to clone should be set to given path:", c.dist)
	}

	os.Setenv("GIT_EXECUTABLE_PATH", "/path/to/git")
	c = newCloner("/path/to/dist")
	if c.git != "/path/to/git" {
		t.Error("Git command should respect environment variable $GIT_EXECUTABLE_PATH:", c.git)
	}

	os.Setenv("GIT_EXECUTABLE_PATH", "")
}

func TestCloneRepo(t *testing.T) {
	c := newCloner("test")
	defer func() {
		os.RemoveAll("test")
	}()
	go c.clone("rhysd/clever-f.vim")
	select {
	case <-c.done:
		s, err := os.Stat("test/rhysd/clever-f.vim")
		if err != nil {
			t.Fatal("Cloned directory not found:", err)
		}
		if !s.IsDir() {
			t.Fatal("It should clone directory")
		}
		if c.err != nil {
			t.Fatal("Error was reported for existing repo:", err)
		}
	}
}
