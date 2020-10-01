package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/github-clone-all/ghca"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const version = "2.4.0"

const usageHeader = `USAGE: github-clone-all [FLAGS] {query}

  github-clone-all is a command to clone all repositories matching to given
  query via GitHub Search API. Query must not be empty.
  It clones many repositories in parallel.

  Repository is cloned to 'dest' directory. It is $cwd/repos by default and
  can be specified with -dest flag.

  All arguments in {query} are regarded as query.
  For example,

  $ github-clone-all foo bar

  will search 'foo bar'. But quoting the query is recommended to avoid
  conflicting with shell special characters as following:

  $ github-clone-all 'foo bar'

  Because of restriction of GitHub search API, max number of results is 1000
  repositories. And you may need to gain GitHub API token in advance to avoid
  reaching API rate limit.

  You can get the token as following:

  1. Visit https://github.com/settings/tokens in a browser
  2. Click 'Generate new token'
  3. Add token description
  4. Without checking any checkbox, click 'Generate token'
  5. Key is shown in your tokens list

  Document for GitHub Repository Search:
    https://help.github.com/articles/searching-repositories/

  Search Syntax:
    https://help.github.com/articles/understanding-the-search-syntax/

  Search API Documentation:
    https://developer.github.com/v3/search/


EXAMPLE:

  $ github-clone-all -extract '(\.vim|vimrc)$' 'language:vim fork:false stars:>1'

    Above command will clone first 1000 repositories into 'repos' directory in
    the current working directory. Only files whose name ending with '.vim' or
    'vimrc' remain in each repositories.

    Query condition:
      - language is 'vim'
      - not a fork repo
      - stars of repo is more than 1

  $ github-clone-all -count 1 'language:javascript'

    Above command will clone the most popular repository of JavaScript on
    GitHub.

  $ github-clone-all -dry 'language:go'

    Above command will only list up most popular 1000 repositories of Go
    instead of cloning them.

  $ github-clone-all -deep -ssh 'user:YOUR_USER_NAME fork:false'

    Above command will clone all your repositories (except for forks) with
    full history. It's useful when you want to clone all your repositories.

FLAGS:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func selfUpdate() int {
	v := semver.MustParse(version)

	latest, err := selfupdate.UpdateSelf(v, "rhysd/github-clone-all")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 3
	}

	if v.Equals(latest.Version) {
		fmt.Println("Current version", v, "is the latest")
	} else {
		fmt.Println("Successfully updated to version", v)
		fmt.Println("Release Note:\n", latest.ReleaseNotes)
	}
	return 0
}

func main() {
	help := flag.Bool("help", false, "Show this help")
	h := flag.Bool("h", false, "Show this help")
	token := flag.String("token", "", "GitHub token to call GitHub API. $GITHUB_TOKEN environment variable is also referred")
	dest := flag.String("dest", "", "Directory to store the downloaded files. By default 'repos' in current working directory")
	extract := flag.String("extract", "", "Regular expression to extract files by name in each cloned repo")
	quiet := flag.Bool("quiet", false, "Run quietly. When exit status is non-zero, it means error occurred")
	count := flag.Int("count", 0, "Max number of repositories to clone")
	dry := flag.Bool("dry", false, "Do dry run. Only shows which repositories will be cloned by given query with repositorie's descriptions")
	deep := flag.Bool("deep", false, "Do not use shallow clone")
	ssh := flag.Bool("ssh", false, "Use git@github.com/... URL instead of https://github.com/... URL")
	ver := flag.Bool("version", false, "Show version")
	update := flag.Bool("selfupdate", false, "Update this tool to the latest")

	flag.Usage = usage
	flag.Parse()

	if *help || *h {
		usage()
		os.Exit(0)
	}

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	if *update {
		os.Exit(selfUpdate())
	}

	if *quiet {
		log.SetOutput(ioutil.Discard)
	}

	query := strings.Join(flag.Args(), " ")

	cli, err := ghca.NewCLI(*token, query, *dest, *extract, *count, *dry, *deep, *ssh)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	if err = cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
