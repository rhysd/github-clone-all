package main

import (
	"os"
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

func TestCollectRepos(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("Skipping because API token not found")
	}
}
