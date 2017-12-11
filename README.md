Clone matching repos on GitHub
==============================

`github-clone-all` is a small command to clone all repositories matching to given query and
language via [GitHub Search API][].
It clones many repositories in parallel. Please see `-help` option to know all flags.

Query is the same as GitHub search syntax. And 'stars>1 fork:false' is added by default for
sensible search results.

Repository is cloned to 'dest' directory. It is $cwd/repos by default and can be specified with
`-dest` flag. And in order to reduce size of cloned repositories, `-extract` option is available.
`-extract` only leaves files matching to given regular expression.

Because of restriction of GitHub search API, max number of results is 1000. And you need to
gain GitHub API token in advance. `github-clone-all` will refer the token via `-token` flag or
`$GITHUB_TOKEN` environment variable.

## Installation

Use `go get` or [released binaries](https://github.com/rhysd/github-clone-all/releases).

```
$ go get github.com/rhysd/github-clone-all
```

## Example

```
$ github-clone-all -token $GITHUB_TOKEN -lang vim -extract '(\.vim|vimrc)$'
```

It clones first 1000 repositories whose language is 'vim' into 'repos' directory in the current
working directory.

## How to get GitHub API token

1. Visit https://github.com/settings/tokens in a browser
2. Click 'Generate new token'
3. Add token description
4. Without checking any checkbox, click 'Generate token'
5. Key is shown in your tokens list

## License

[MIT license](LICENSE)

[GitHub Search API]: https://developer.github.com/v3/search/
