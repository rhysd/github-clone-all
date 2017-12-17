Clone matching repos on GitHub
==============================
[![GoDoc Badge][]][GoDoc]
[![Mac and Linux Build Status][]][Travis CI]
[![Windows Build Status][]][Appveyor]
[![Coverage Status][]][Codecov]

```
$ github-clone-all [flags] {query}
```

`github-clone-all` is a small command to clone all repositories matching to given query and
language via [GitHub Search API][]. The query should be in [GitHub search syntax][] and cannot be
empty. It clones many repositories in parallel. Please see `-help` option to know all flags.

Repository is cloned to 'dest' directory. It is `$cwd/repos` by default and can be specified with
`-dest` flag. And in order to reduce size of cloned repositories, `-extract` option is available.
`-extract` only leaves files matching to given regular expression.

Because of restriction of GitHub search API, max number of results is 1000 repositories. And you may
need to gain GitHub API token in advance to avoid reaching API rate limit. `github-clone-all` will
refer the token via `-token` flag or `$GITHUB_TOKEN` environment variable.

All arguments in `{query}` are regarded as query. For example, `github-clone-all foo bar` will search
`foo bar`. But quoting the query is recommended to avoid conflicting with shell special characters
as `github-clone-all 'foo bar'`.


## Installation

Use `go get` or [released binaries](https://github.com/rhysd/github-clone-all/releases).

```
$ go get github.com/rhysd/github-clone-all
```


## Example

```
$ github-clone-all -extract '(\.vim|vimrc)$' 'language:vim fork:false stars:>1'
```

Above command will clone first 1000 repositories into 'repos' directory in the current working
directory. And it only leaves files whose file name ends with `.vim` or `vimrc`.
So it collects many Vim script codes from famous repositories on GitHub.

Query condition:
- language is 'vim'
- not a fork repo
- stars of repo is more than 1

```
$ github-clone-all -count 1 'language:javascript'
```

Above command will clone the most popular repository of JavaScript on GitHub.

```
$ github-clone-all -dry 'language:go'
```

Above command will only list up most popular 1000 repositories of Go instead of cloning them.


## How to get GitHub API token

1. Visit https://github.com/settings/tokens in a browser
2. Click 'Generate new token'
3. Add token description
4. Without checking any checkbox, click 'Generate token'
5. Generated token is shown at the top of your tokens list


## Use github-clone-all programmatically

`github-clone-all` consists of tiny `main.go` and `ghca` package. You can import `ghca` to utilize
functions of the tool.

```go
import "github.com/rhysd/github-clone-all/ghca"
```

Please read [documentation][GoDoc] for more detail.

## License

[MIT license](LICENSE)

[GitHub search syntax]: https://help.github.com/articles/understanding-the-search-syntax/
[GitHub Search API]: https://developer.github.com/v3/search/
[GoDoc Badge]: https://godoc.org/github.com/rhysd/github-clone-all/ghca?status.svg
[GoDoc]: https://godoc.org/github.com/rhysd/github-clone-all/ghca
[Mac and Linux Build Status]: https://travis-ci.org/rhysd/github-clone-all.svg?branch=master
[Travis CI]: https://travis-ci.org/rhysd/github-clone-all
[Windows Build Status]: https://ci.appveyor.com/api/projects/status/fwaaouneyn9kftts/branch/master?svg=true
[Appveyor]: https://ci.appveyor.com/project/rhysd/github-clone-all/branch/master
[Coverage Status]: https://codecov.io/gh/rhysd/github-clone-all/branch/master/graph/badge.svg
[Codecov]: https://codecov.io/gh/rhysd/github-clone-all
