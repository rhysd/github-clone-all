package main

import (
	"flag"
	"fmt"
	"github.com/rhysd/github-clone-all/ghca"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const usageHeader = `USAGE: github-clone-all [FLAGS] {query}

  github-clone-all is a command to clone all repositories matching to given
  query via GitHub Search API. Query must not be empty.
  It clones many repositories in parallel.

  Repository is cloned to 'dest' directory. It is $cwd/repos by default and
  can be specified with -dest flag.

  Because of restriction of GitHub search API, max number of results is 1000.
  And you need to gain GitHub API token in advance to avoid API rate limit.

  All arguments in {query} are regarded as query.
  For example,

  $ github-clone-all foo bar

  will search 'foo bar'. But quoting the query is recommended to avoid
  conflicting with shell special characters as following:

  $ github-clone-all 'foo bar'

  You can get the token as following:

  1. Visit https://github.com/settings/tokens in a browser
  2. Click 'Generate new token'
  3. Add token description
  4. Without checking any checkbox, click 'Generate token'
  5. Key is shown in your tokens list

  ref: https://developer.github.com/v3/search/


EXAMPLE:

  $ github-clone-all -token xxxxxxxx -extract '(\.vim|vimrc)$' 'language:vim fork:false stars:>1'

  It clones first 1000 repositories into 'repos' directory in the current
  working directory.

  Query condition:
    - language is 'vim'
    - not a fork repo
    - stars of repo is more than 1

  If the token is set to $GITHUB_TOKEN environment variable, following should
  also work fine.

  $ github-clone-all -extract '(\.vim|vimrc)$' 'language:vim fork:false stars:>1'


FLAGS:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	help := flag.Bool("help", false, "Show this help")
	h := flag.Bool("h", false, "Show this help")
	token := flag.String("token", "", "GitHub token to call GitHub API. If this option is not specified, $GITHUB_TOKEN environment variable needs to be set")
	dest := flag.String("dest", "", "Directory to store the downloaded files. By default 'repos' in current working directory")
	extract := flag.String("extract", "", "Regular expression to extract files in each cloned repo")
	quiet := flag.Bool("quiet", false, "Run quietly. Exit status is non-zero, it means error occurred")

	flag.Usage = usage
	flag.Parse()

	if *help || *h {
		usage()
		os.Exit(0)
	}

	if *quiet {
		log.SetOutput(ioutil.Discard)
	}

	query := strings.Join(flag.Args(), " ")

	cli, err := ghca.NewCLI(*token, query, *dest, *extract)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	if err = cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
