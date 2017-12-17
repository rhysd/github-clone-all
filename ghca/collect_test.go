package ghca

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector("foo", "", "", nil, 0, false, nil)
	if c.perPage != 100 {
		t.Error("perPage should be 100 by default:", c.perPage)
	}
	if c.maxPage != 10 {
		t.Error("maxPage should be 10 by default:", c.maxPage)
	}
	if c.page != 1 {
		t.Error("page should be 1 by default:", c.page)
	}
	if c.Query != "foo" {
		t.Error("query should be set to 'foo'", c.Query)
	}
}

func TestNewCollectorWithConfig(t *testing.T) {
	c := NewCollector("foo", "", "", nil, 0, false, &PageConfig{1, 10, 3})
	if c.perPage != 1 {
		t.Error("perPage should be set to 1:", c.perPage)
	}
	if c.maxPage != 10 {
		t.Error("maxPage should be set to 10:", c.maxPage)
	}
	if c.page != 3 {
		t.Error("page should be set to 3:", c.page)
	}

	c = NewCollector("foo", "", "", nil, 0, false, &PageConfig{3, PageUnlimited, 3})
	if c.maxPage != 334 {
		t.Error("maxPage should be calculated to fetch 1000 repos:", c.maxPage)
	}
}

func TestCollectReposTotalIsAFew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping because API token not found")
	}

	defer func() {
		os.RemoveAll("test")
	}()

	c := NewCollector("clever-f.vim language:vim fork:false", token, "test", nil, 0, false, nil)
	count, total, err := c.Collect()
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
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping because API token not found")
	}

	defer func() {
		os.RemoveAll("test")
	}()

	// Get page 4, 5, 6 and each page results in 2 repos
	c := NewCollector("language:vim fork:false", token, "test", nil, 0, false, &PageConfig{
		Per:   2,
		Max:   6,
		Start: 4,
	})

	count, total, err := c.Collect()
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

func TestBadCredential(t *testing.T) {
	defer func() {
		os.RemoveAll("test")
	}()
	c := NewCollector("clever-f.vim language:vim fork:false", "badcredentials", "test", nil, 0, false, nil)
	_, _, err := c.Collect()
	if err == nil {
		t.Fatal("Bad credentials should cause an error on collecting")
	}
}

func TestSpecifyCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	defer func() {
		os.RemoveAll("test")
	}()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping because API token not found")
	}

	c := NewCollector("user:rhysd", token, "test", nil, 2, false, nil)
	if c.maxPage != 1 {
		t.Fatal("Max page should be 1 if count is specified as 2 because of 100 repos per page:", c.maxPage)
	}

	count, _, err := c.Collect()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal("Count is specified as 2 but actually 2 repos are not cloned:", count)
	}

	fs, err := ioutil.ReadDir(filepath.FromSlash("test/rhysd"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 2 {
		ns := make([]string, 0, len(fs))
		for _, f := range fs {
			ns = append(ns, f.Name())
		}
		t.Fatal("Count is specified as 2 but actually 2 repos are not cloned:", count, ", ", strings.Join(ns, " "))
	}
}

func TestDryRun(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping because API token not found")
	}

	c := NewCollector("user:rhysd", token, "test", nil, 2, true, nil)
	_, _, err := c.Collect()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat("test"); err == nil {
		os.RemoveAll("test")
		t.Fatal("'test' directory was created in spite of dry-run")
	}
}
