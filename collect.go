package main

type collector struct {
	perPage uint
	maxPage uint
	page    uint
	query   string
}

type pageConfig struct {
	per   uint
	max   uint
	start uint
}

const pageUnlimited uint = 0

func newCollector(query string, page *pageConfig) *collector {
	c := &collector{100, pageUnlimited, 1, query}
	if page != nil {
		c.perPage = page.per
		c.maxPage = page.max
		c.page = page.start
	}
	return c
}
