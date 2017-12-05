package main

import (
	"os"
	"testing"
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
