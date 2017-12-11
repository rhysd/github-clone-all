package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCollector(t *testing.T) {
	c := newCollector("foo", "", "", nil, nil)
	if c.perPage != 100 {
		t.Error("perPage should be 100 by default:", c.perPage)
	}
	if c.maxPage != 10 {
		t.Error("maxPage should be 10 by default:", c.maxPage)
	}
	if c.page != 1 {
		t.Error("page should be 1 by default:", c.page)
	}
	if c.query != "foo" {
		t.Error("query should be set to 'foo'", c.query)
	}
}

func TestNewCollectorWithConfig(t *testing.T) {
	c := newCollector("foo", "", "", nil, &pageConfig{1, 10, 3})
	if c.perPage != 1 {
		t.Error("perPage should be set to 1:", c.perPage)
	}
	if c.maxPage != 10 {
		t.Error("maxPage should be set to 10:", c.maxPage)
	}
	if c.page != 3 {
		t.Error("page should be set to 3:", c.page)
	}

	c = newCollector("foo", "", "", nil, &pageConfig{3, pageUnlimited, 3})
	if c.maxPage != 334 {
		t.Error("maxPage should be calculated to fetch 1000 repos:", c.maxPage)
	}
}

func TestCollectReposTotalIsAFew(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping because API token not found")
	}

	defer func() {
		os.RemoveAll("test")
	}()

	c := newCollector("clever-f.vim language:vim fork:false", token, "test", nil, nil)
	count, total, err := c.collect()
	if err != nil {
		t.Fatal("Failed to collect", err)
	}
	if total < 2 || count < 2 {
		t.Fatal("Total repositories is too few:", total)
	}

	for _, dir := range []string{
		"test/rhysd/clever-f.vim",
		"test/vim-scripts/clever-f.vim",
	} {
		dir = filepath.FromSlash(dir)
		if _, err := os.Stat(dir); err != nil {
			t.Fatal(dir, "was not cloned")
		}
	}
}

func TestCollectReposTotalIsLarge(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping because API token not found")
	}

	defer func() {
		os.RemoveAll("test")
	}()

	// Get page 4, 5, 6 and each page results in 2 repos
	c := newCollector("language:vim fork:false", token, "test", nil, &pageConfig{
		per:   2,
		max:   6,
		start: 4,
	})

	count, total, err := c.collect()
	if err != nil {
		t.Fatal("Failed to collect", err)
	}
	if total == 0 {
		t.Fatal("No repository was found")
	}
	if count != 6 {
		t.Fatal("6 repositories (2x3) should be resulted:", total)
	}

	dir, err := os.Open("test")
	if err != nil {
		t.Fatal("Failed to open", err)
	}
	if _, err := dir.Readdirnames(1); err == io.EOF {
		t.Fatal("'test' directory is empty")
	}
}
