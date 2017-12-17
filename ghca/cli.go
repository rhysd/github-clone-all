// Package ghca provides functionalities of github-clone-all command.
// Because of restriction of GitHub search API, max number of results is 1000 repositories.
// And you may need to gain GitHub API token in advance to avoid reaching API rate limit.
//
// Please see the repository page to know more detail.
//   https://github.com/rhysd/github-clone-all
package ghca

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CLI represents a command line interface of github-clone-all.
type CLI struct {
	token   string
	query   string
	dest    string
	extract *regexp.Regexp // Maybe nil
	count   int
}

func (c *CLI) ensureReposDir() error {
	s, err := os.Stat(c.dest)
	if err != nil {
		return os.Mkdir(c.dest, 0755)
	}
	if !s.IsDir() {
		return fmt.Errorf("Cannot create directory '%s' because it's a file", c.dest)
	}
	return nil
}

// Run processes github-clone-all with given options.
func (c *CLI) Run() (err error) {
	if err = c.ensureReposDir(); err != nil {
		return
	}
	col := NewCollector(c.query, c.token, c.dest, c.extract, c.count, nil)
	_, _, err = col.Collect()
	return
}

// NewCLI creates a new command line interface to run github-clone-all.
// Query ('q' parameter) must not be empty.
func NewCLI(t, q, d, e string, c int) (*CLI, error) {
	var err error

	if t == "" {
		t = os.Getenv("GITHUB_TOKEN")
	}

	if d == "" {
		d, err = os.Getwd()
		if err != nil {
			return nil, err
		}
		d = filepath.Join(d, "repos")
	}

	var r *regexp.Regexp
	if e != "" {
		r, err = regexp.Compile(e)
		if err != nil {
			return nil, err
		}
	}

	q = strings.TrimSpace(q)
	if q == "" {
		return nil, errors.New("Query cannot be empty")
	}

	return &CLI{t, q, d, r, c}, nil
}
