package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help  = flag.Bool("help", false, "Show this help")
	token = flag.String("token", "", "GitHub token to call GitHub API")
	regex = flag.String("regex", "", "Regular expression matching file names to downoad")
	lang  = flag.String("lang", "", "Language name to search repos")
	dist  = flag.String("dist", "", "Directory to store the downloaded files. Current working directory by default")
)

const usageHeader = `Usage: repo-collect-gh -token {token} -regex {regex} -lang {lang} [-dist {path}]

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

	os.Exit(start(*token, *regex, *lang, *dist))
}
