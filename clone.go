package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
)

const maxConcurrency = 4
const maxBuffer = 1000

type cloner struct {
	git     string
	dist    string
	extract *regexp.Regexp
	repos   chan string
	err     chan error
	wg      sync.WaitGroup
}

func newCloner(dist string, extract *regexp.Regexp) *cloner {
	c := &cloner{
		git:     os.Getenv("GIT_EXECUTABLE_PATH"),
		dist:    dist,
		extract: extract,
		repos:   make(chan string, maxBuffer),
		err:     make(chan error),
	}

	if c.git == "" {
		c.git = "git"
	}

	return c
}

func (cl *cloner) clone(repo string) {
	cl.repos <- repo
}

func (cl *cloner) newWorker() {
	cl.wg.Add(1)
	dist := cl.dist
	git := cl.git

	var extract *regexp.Regexp
	if cl.extract != nil {
		extract = cl.extract.Copy()
	}

	go func() {
		defer cl.wg.Done()
		for repo := range cl.repos {
			log.Println("Cloning", repo)

			url := fmt.Sprintf("https://github.com/%s.git", repo)
			dir := filepath.FromSlash(fmt.Sprintf("%s/%s", dist, repo))
			cmd := exec.Command(git, "clone", "--depth=1", "--single-branch", url, dir)
			err := cmd.Run()

			if err != nil {
				log.Println("Failed to clone", repo, err)
				cl.err <- err
				return
			}

			if extract != nil {
				if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}
					if extract.MatchString(path) {
						if err := os.Remove(path); err != nil {
							return err
						}
					}
					return nil
				}); err != nil {
					log.Println("Failed to extract files", repo, extract.String(), err)
					cl.err <- err
					return
				}
			}

			log.Println("Cloned:", repo)
		}
	}()
}

func (cl *cloner) start() {
	for i := 0; i < maxConcurrency; i++ {
		cl.newWorker()
	}
}

func (cl *cloner) shutdown() {
	close(cl.repos)
	cl.wg.Wait()
}
