package ghca

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli, err := NewCLI("token", "foo stars>1", "dest", "")
	if err != nil {
		t.Fatal(err)
	}
	if cli.token != "token" {
		t.Error("Unexpected token", cli.token)
	}
	if cli.query != "foo stars>1" {
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
	cli, err := NewCLI("token", "query", "", "")
	if err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	d := filepath.Join(cwd, "repos")
	if cli.dest != d {
		t.Error("Empty dest should mean current working directory but:", cli.dest)
	}
}

func TestEmptyToken(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "")
	if _, err := NewCLI("", "query", "", ""); err == nil {
		t.Error("Empty token should raise an error")
	}
	os.Setenv("GITHUB_TOKEN", token)
}

func TestEmptyQuery(t *testing.T) {
	for _, q := range []string{
		"",
		"   ",
		"	",
	} {
		if _, err := NewCLI("token", q, "", ""); err == nil {
			t.Errorf("Empty query should raise an error: '%s'", q)
		}
	}
}

func TestGitHubTokenEnv(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "foobar")
	cli, err := NewCLI("", "query", "", "")
	if err != nil {
		t.Error(err)
	}
	if cli.token != "foobar" {
		t.Error("Unexpected token", cli.token)
	}
	os.Setenv("GITHUB_TOKEN", token)
}

func TestInvalidRegexp(t *testing.T) {
	if _, err := NewCLI("token", "query", "", "(foo"); err == nil {
		t.Error("Broken regexp must raise an error")
	}
}

func TestMakeDest(t *testing.T) {
	defer os.Remove("repos")

	cli, err := NewCLI("token", "query", "", "")
	if err != nil {
		t.Fatal(err)
	}

	// If directory is already existing, it does nothing. First create the directory
	// and at second check the case where directory is already existing.
	for i := 0; i < 2; i++ {
		if err := cli.ensureReposDir(); err != nil {
			t.Fatal(err)
		}

		s, err := os.Stat("repos")
		if err != nil {
			t.Fatal("'repos' directory was not created")
		}
		if !s.IsDir() {
			t.Fatal("Created entry is not a directory")
		}
	}
}

func TestDestAlreadyExistAsFile(t *testing.T) {
	defer os.Remove("repos")
	f, err := os.OpenFile("repos", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	cli, err := NewCLI("token", "query", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := cli.ensureReposDir(); err == nil {
		t.Fatal("Error should occur when file is already created")
	}
}
