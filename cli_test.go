package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli, err := newCLI("token", "foo stars>1", "lang", "dest", "")
	if err != nil {
		t.Fatal(err)
	}
	if cli.token != "token" {
		t.Error("Unexpected token", cli.token)
	}
	if cli.query != "foo stars>1 language:lang fork:false" {
		t.Error("Unexpected query", cli.query)
	}
	if cli.dest != "dest" {
		t.Error("Unexpected dest", cli.dest)
	}
	if cli.extract != nil {
		t.Error("Invalid regular expression for empty extract pattern:", *cli.extract)
	}
}

func TestEmptyDest(t *testing.T) {
	cli, err := newCLI("token", "query", "lang", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	d := filepath.Join(cwd, "repos")
	if cli.dest != d {
		t.Error("Empty dest should mean current working directory but:", cli.dest)
	}
}

func TestEmptyTokenOrLang(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "")
	if _, err := newCLI("", "", "vim", "", ""); err == nil {
		t.Error("Empty token should raise an error")
	}

	if _, err := newCLI("", "foobar", "", "", ""); err == nil {
		t.Error("Empty lang should raise an error")
	}
	os.Setenv("GITHUB_TOKEN", token)
}

func TestGitHubTokenEnv(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "foobar")
	cli, err := newCLI("", "", "vim", "", "")
	if err != nil {
		t.Error(err)
	}
	if cli.token != "foobar" {
		t.Error("Unexpected token", cli.token)
	}
	os.Setenv("GITHUB_TOKEN", token)
}

func TestInvalidRegexp(t *testing.T) {
	if _, err := newCLI("token", "", "vim", "", "(foo"); err == nil {
		t.Error("Broken regexp must raise an error")
	}

}
