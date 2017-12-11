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
	dest    string
	extract *regexp.Regexp
	repos   chan string
	err     chan error
	wg      sync.WaitGroup
	ssh     bool
}

func newCloner(dest string, extract *regexp.Regexp) *cloner {
	c := &cloner{
		git:     os.Getenv("GIT_EXECUTABLE_PATH"),
		dest:    dest,
		extract: extract,
		repos:   make(chan string, maxBuffer),
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
	dest := cl.dest
	git := cl.git

	var extract *regexp.Regexp
	if cl.extract != nil {
		extract = cl.extract.Copy()
	}

	go func() {
		defer cl.wg.Done()
		for repo := range cl.repos {
			log.Println("Cloning", repo)

			var url string
			if cl.ssh {
				url = fmt.Sprintf("git@github.com:%s.git", repo)
			} else {
				url = fmt.Sprintf("https://github.com/%s.git", repo)
			}

			dir := filepath.FromSlash(fmt.Sprintf("%s/%s", dest, repo))
			cmd := exec.Command(git, "clone", "--depth=1", "--single-branch", url, dir)
			err := cmd.Run()

			if err != nil {
				log.Println("Failed to clone", repo, err)
				if cl.err != nil {
					cl.err <- err
				}
				continue
			}

			if extract != nil {
				if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}
					if !extract.MatchString(path) {
						if err := os.Remove(path); err != nil {
							return err
						}
					}
					return nil
				}); err != nil {
					log.Println("Failed to extract files", repo, extract.String(), err)
					if cl.err != nil {
						cl.err <- err
					}
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
	if cl.err != nil {
		close(cl.err)
	}
}
