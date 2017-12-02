package main

import (
	"fmt"
	"os"
	"regexp"
)

type cli struct {
	token string
	regex *regexp.Regexp
	lang  string
	dist  string
}

func (c *cli) run() (err error) {
	return
}

func newCLI(t, r, l, d string) (*cli, error) {
	if t == "" || r == "" || l == "" {
		return nil, fmt.Errorf("All of token, regex and lang must be set. Please see -help for more detail")
	}

	re, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}

	if d == "" {
		d, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	return &cli{t, re, l, d}, nil
}

func start(token, regex, lang, dist string) int {
	c, err := newCLI(token, regex, lang, dist)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 3
	}
	if err = c.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 3
	}
	return 0
}
