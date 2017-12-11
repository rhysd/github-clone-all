package ghca

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

type Collector struct {
	perPage uint
	maxPage uint
	page    uint
	Query   string
	Dest    string
	Extract *regexp.Regexp
	client  *github.Client
	ctx     context.Context
}

func (col *Collector) searchRepos() (*github.RepositoriesSearchResult, error) {
	o := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    int(col.page),
			PerPage: int(col.perPage),
		},
	}
	r, _, err := col.client.Search.Repositories(col.ctx, col.Query, o)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (col *Collector) Collect() (int, int, error) {
	log.Println("Searching GitHub repositories with query:", col.Query)
	cloner := NewCloner(col.Dest, col.Extract)
	cloner.Start()

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
			cloner.Clone(fmt.Sprintf("%s/%s", repo.GetOwner().GetLogin(), repo.GetName()))
			count++
		}

		col.page++
	}

	cloner.Shutdown()

	log.Println(count, "repositories were cloned into", col.Dest, "for total", total, "search results")

	return count, total, nil
}

type PageConfig struct {
	Per   uint
	Max   uint
	Start uint
}

const PageUnlimited uint = 0

func NewCollector(query, token, dest string, extract *regexp.Regexp, page *PageConfig) *Collector {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(ctx, src))
	c := &Collector{100, PageUnlimited, 1, query, dest, extract, client, ctx}
	if page != nil {
		c.perPage = page.Per
		c.maxPage = page.Max
		c.page = page.Start
	}
	if c.maxPage == PageUnlimited {
		c.maxPage = uint(math.Ceil(1000.0 / float64(c.perPage)))
	}
	return c
}
