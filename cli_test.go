package main

import (
	"os"
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli, err := newCLI("token", "query", "lang", "dist")
	if err != nil {
		t.Fatal(err)
	}
	if cli.token != "token" {
		t.Error("Unexpected token", cli.token)
	}
	if cli.query != "query" {
		t.Error("Unexpected query", cli.query)
	}
	if cli.lang != "lang" {
		t.Error("Unexpected lang", cli.lang)
	}
	if cli.dist != "dist" {
		t.Error("Unexpected dist", cli.dist)
	}
}

func TestEmptyDist(t *testing.T) {
	cli, err := newCLI("token", "query", "lang", "")
	if err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	if cli.dist != cwd {
		t.Error("Empty dist should mean current working directory but:", cli.dist)
	}
}

func TestEmptyTokenOrLang(t *testing.T) {
	if _, err := newCLI("", "", "vim", ""); err == nil {
		t.Error("Empty token should raise an error")
	}

	os.Setenv("GITHUB_TOKEN", "")
	if _, err := newCLI("", "foobar", "", ""); err == nil {
		t.Error("Empty lang should raise an error")
	}
}

func TestGitHubTokenEnv(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "foobar")
	cli, err := newCLI("", "", "vim", "")
	if err != nil {
		t.Fatal(err)
	}
	if cli.token != "foobar" {
		t.Fatal("Unexpected token", cli.token)
	}
}
