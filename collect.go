package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"math"
	"regexp"
	"time"
)

type collector struct {
	perPage uint
	maxPage uint
	page    uint
	query   string
	dest    string
	extract *regexp.Regexp
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

func (col *collector) collect() (int, int, error) {
	log.Println("Searching GitHub repositories with query:", col.query)
	cloner := newCloner(col.dest, col.extract)
	cloner.start()

	total := 0
	count := 0
	for col.page <= col.maxPage {
		res, err := col.searchRepos()
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("Rate limit exceeded. Sleeping 1 minute")
			time.Sleep(1 * time.Minute)
			continue
		} else if err != nil {
			return 0, 0, err
		}

		total = res.GetTotal()

		if res.GetIncompleteResults() {
			log.Println("TODO: Handle incomplete result returned from GitHub API")
		}

		if len(res.Repositories) == 0 {
			// All repositories were searched
			break
		}

		for _, repo := range res.Repositories {
			cloner.clone(fmt.Sprintf("%s/%s", repo.GetOwner().GetLogin(), repo.GetName()))
			count++
		}

		col.page++
	}

	cloner.shutdown()

	log.Println(count, "repositories were cloned into", col.dest, "for total", total, "search results")

	select {
	case err, ok := <-cloner.err:
		if ok {
			return 0, 0, err
		}
	default:
		// Do nothing
	}

	return count, total, nil
}

type pageConfig struct {
	per   uint
	max   uint
	start uint
}

const pageUnlimited uint = 0

func newCollector(query, token, dest string, extract *regexp.Regexp, page *pageConfig) *collector {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(ctx, src))
	c := &collector{100, pageUnlimited, 1, query, dest, extract, client, ctx}
	if page != nil {
		c.perPage = page.per
		c.maxPage = page.max
		c.page = page.start
	}
	if c.maxPage == pageUnlimited {
		c.maxPage = uint(math.Ceil(1000.0 / float64(c.perPage)))
	}
	return c
}
