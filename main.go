package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help  = flag.Bool("help", false, "Show this help")
	token = flag.String("token", "", "GitHub token to call GitHub API")
	query = flag.String("query", "", "Additional query string to search")
	lang  = flag.String("lang", "", "Language name to search repos")
	dist  = flag.String("dist", "", "Directory to store the downloaded files. Current working directory by default")
)

const usageHeader = `Usage: repo-collect-gh -token {token} -lang {lang} [-query {query}] [-dist {path}]

  Under construction!
  Description goes here.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	cli, err := newCLI(*token, *query, *lang, *dist)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	if err = cli.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
