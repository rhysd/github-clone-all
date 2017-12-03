package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"math"
)

type collector struct {
	perPage uint
	maxPage uint
	page    uint
	query   string
	client  *github.Client
	ctx     context.Context
}

func (col *collector) searchRepos() (*github.RepositoriesSearchResult, error) {
	o := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    int(col.page),
			PerPage: int(col.perPage),
		},
	}
	r, _, err := col.client.Search.Repositories(col.ctx, col.query, o)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (col *collector) collect() error {
	return nil // TODO
}

type pageConfig struct {
	per   uint
	max   uint
	start uint
}

const pageUnlimited uint = 0

func newCollector(query string, token string, page *pageConfig) *collector {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(ctx, src))
	c := &collector{100, pageUnlimited, 1, query, client, ctx}
	if page != nil {
		c.perPage = page.per
		c.maxPage = page.max
		c.page = page.start
	}
	if c.maxPage == 0 {
		c.maxPage = uint(math.Ceil(1000.0 / float64(c.perPage)))
	}
	return c
}
