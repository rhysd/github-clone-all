package main

import (
	"flag"
	"fmt"
	"os"
)

const usageHeader = `Usage: github-clone-all {Flags}

  github-clone-all is a command to clone all repositories matching to given
  query and language via GitHub Search API.
  It clones many repositories in parallel.

  Query is the same as GitHub search syntax. And 'stars>1 fork:false' is
  added by default for sensible search results.

  Repository is cloned to 'dest' directory. It is $cwd/repos by default and
  can be specified with -dest flag.

  Because of restriction of GitHub search API, max number of results is 1000.
  And you need to gain GitHub API token in advance. You can get the token as
  following:

  1. Visit https://github.com/settings/tokens in a browser
  2. Click 'Generate new token'
  3. Add token description
  4. Without checking any checkbox, click 'Generate token'
  5. Key is shown in your tokens list

  ref: https://developer.github.com/v3/search/

Example:

  $ github-clone-all -token $GITHUB_TOKEN -lang vim -extract '(\.vim|vimrc)$'

  It clones first 1000 repositories whose language is 'vim' into 'repos'
  directory in the current working directory.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	help := flag.Bool("help", false, "Show this help")
	h := flag.Bool("h", false, "Show this help")
	token := flag.String("token", "", "GitHub token to call GitHub API. $GITHUB_TOKEN environment variable is also referred (required)")
	query := flag.String("query", "", "Additional query string to search (optional)")
	lang := flag.String("lang", "", "Language name to search repos (required)")
	dest := flag.String("dest", "", "Directory to store the downloaded files. By default 'repos' in current working directory (optional)")
	extract := flag.String("extract", "", "Regular expression to extract files in each cloned repo (optional)")

	flag.Usage = usage
	flag.Parse()

	if *help || *h {
		usage()
		os.Exit(0)
	}

	cli, err := newCLI(*token, *query, *lang, *dest, *extract)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	if err = cli.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
