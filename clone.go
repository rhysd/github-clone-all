package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type worker struct {
	git   string
	dist  string
	repos chan string
	ret   chan error
	wg    *sync.WaitGroup
}

func (w *worker) start() {
	defer w.wg.Done()
	for repo := range w.repos {
		log.Println("Cloning", repo)
		url := fmt.Sprintf("https://github.com/%s.git", repo)
		dir := fmt.Sprintf("%s/%s", w.dist, repo)
		cmd := exec.Command(w.git, "clone", "--depth=1", "--single-branch", url, dir)
		err := cmd.Run()
		if err != nil {
			w.ret <- err
		}
		log.Println("Cloned:", repo)
	}
}

const maxConcurrency = 4
const maxBuffer = 1000

type cloner struct {
	git   string
	dist  string
	repos chan string
	err   chan error
	wg    sync.WaitGroup
}

func newCloner(dist string) *cloner {
	c := &cloner{
		git:   os.Getenv("GIT_EXECUTABLE_PATH"),
		dist:  dist,
		repos: make(chan string, maxBuffer),
		err:   make(chan error),
	}

	if c.git == "" {
		c.git = "git"
	}

	return c
}

func (cl *cloner) clone(repo string) {
	cl.repos <- repo
}

func (cl *cloner) start() {
	for i := 0; i < maxConcurrency; i++ {
		w := &worker{cl.git, cl.dist, cl.repos, cl.err, &cl.wg}
		cl.wg.Add(1)
		go w.start()
	}
}

func (cl *cloner) shutdown() {
	close(cl.repos)
	cl.wg.Wait()
}
