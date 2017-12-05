package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"math"
	"time"
)

type collector struct {
	perPage uint
	maxPage uint
	page    uint
	query   string
	dist    string
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
	cloners := make([]*cloner, 0, 4)
	for i := 0; i < 4; i++ {
		c := newCloner(col.dist)
		cloners = append(cloners, c)
		go c.start()
	}

	for col.page <= col.maxPage {
		res, err := col.searchRepos()
		if _, ok := err.(*github.RateLimitError); ok {
			time.Sleep(1 * time.Minute)
			continue
		} else if err != nil {
			return err
		}
		if res.GetIncompleteResults() {
			log.Println("TODO: Handle incomplete result returned from GitHub API")
		}

		// TODO: 空いているやつに優先的に割り当てていくスケジューラをつくる（もしくはライブラリを調べて使う）
		for i, repo := range res.Repositories {
			c := cloners[i%4]
			c.recv <- fmt.Sprintf("%s/%s", repo.GetName(), repo.GetOwner().GetLogin())
		}

		col.page++
	}

	for _, c := range cloners {
		close(c.recv) // Sends "" to c.recv
	}

	return nil
}

type pageConfig struct {
	per   uint
	max   uint
	start uint
}

const pageUnlimited uint = 0

func newCollector(query, token, dist string, page *pageConfig) *collector {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(ctx, src))
	c := &collector{100, pageUnlimited, 1, query, dist, client, ctx}
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
