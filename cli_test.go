package main

import (
	"os"
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli, err := newCLI("token", "foo stars>1", "lang", "dist", "")
	if err != nil {
		t.Fatal(err)
	}
	if cli.token != "token" {
		t.Error("Unexpected token", cli.token)
	}
	if cli.query != "foo stars>1 language:lang fork:false" {
		t.Error("Unexpected query", cli.query)
	}
	if cli.dist != "dist" {
		t.Error("Unexpected dist", cli.dist)
	}
	if cli.extract != nil {
		t.Error("Invalid regular expression for empty extract pattern:", *cli.extract)
	}
}

func TestEmptyDist(t *testing.T) {
	cli, err := newCLI("token", "query", "lang", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	if cli.dist != cwd {
		t.Error("Empty dist should mean current working directory but:", cli.dist)
	}
}

func TestEmptyTokenOrLang(t *testing.T) {
	if _, err := newCLI("", "", "vim", "", ""); err == nil {
		t.Error("Empty token should raise an error")
	}

	os.Setenv("GITHUB_TOKEN", "")
	if _, err := newCLI("", "foobar", "", "", ""); err == nil {
		t.Error("Empty lang should raise an error")
	}
}

func TestGitHubTokenEnv(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "foobar")
	cli, err := newCLI("", "", "vim", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if cli.token != "foobar" {
		t.Fatal("Unexpected token", cli.token)
	}
}

func TestInvalidRegexp(t *testing.T) {
	if _, err := newCLI("token", "", "vim", "", "(foo"); err == nil {
		t.Error("Broken regexp must raise an error")
	}

}
