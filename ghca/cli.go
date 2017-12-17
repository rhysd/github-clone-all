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
	dry     bool
}

func (c *CLI) ensureReposDir() error {
	if c.dry {
		return nil
	}
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
	col := NewCollector(c.query, c.token, c.dest, c.extract, c.count, c.dry, nil)
	_, _, err = col.Collect()
	return
}

// NewCLI creates a new command line interface to run github-clone-all.
// Query ('q' parameter) must not be empty.
func NewCLI(token, query, dest, extract string, count int, dry bool) (*CLI, error) {
	var err error

	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	if dest == "" {
		dest, err = os.Getwd()
		if err != nil {
			return nil, err
		}
		dest = filepath.Join(dest, "repos")
	}

	var r *regexp.Regexp
	if extract != "" {
		r, err = regexp.Compile(extract)
		if err != nil {
			return nil, err
		}
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("Query cannot be empty")
	}

	return &CLI{token, query, dest, r, count, dry}, nil
}
