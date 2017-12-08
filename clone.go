package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type cloner struct {
	git     string
	dist    string
	ret     []chan error
	done    []bool
	lastErr error
}

func newCloner(dist string) *cloner {
	git := os.Getenv("GIT_EXECUTABLE_PATH")
	if git == "" {
		git = "git"
	}

	ret := make([]chan error, 0, 4)
	for i := 0; i < 4; i++ {
		ret = append(ret, make(chan error))
	}

	done := make([]bool, 0, 4)
	for i := 0; i < 4; i++ {
		done = append(done, true)
	}

	return &cloner{git, dist, ret, done, nil}
}

// Note: This function is run in another goroutine. It should not share the state with cloner so it should not be a method of cloner.
func clone(git, repo, dist string, done chan error) {
	log.Println("Cloning", repo)
	url := fmt.Sprintf("https://github.com/%s.git", repo)
	dir := fmt.Sprintf("%s/%s", dist, repo)
	cmd := exec.Command(git, "clone", "--depth=1", "--single-branch", url, dir)
	err := cmd.Run()
	log.Println("Cloned:", repo)
	done <- err
}

func (cl *cloner) waitOne() (idx int) {
	var err error

	select {
	case err = <-cl.ret[0]:
		cl.done[0] = true
		idx = 0
	case err = <-cl.ret[1]:
		cl.done[1] = true
		idx = 1
	case err = <-cl.ret[2]:
		cl.done[2] = true
		idx = 2
	case err = <-cl.ret[3]:
		cl.done[3] = true
		idx = 3
	}

	if err != nil {
		log.Println("Failed to clone:", err)
		cl.lastErr = err
	}
	return
}

func (cl *cloner) waitDone() {
	for !cl.done[0] || !cl.done[1] || !cl.done[2] || !cl.done[3] {
		cl.waitOne()
	}
}

func (cl *cloner) cloneWith(idx int, repo string) {
	cl.done[idx] = false
	go clone(cl.git, repo, cl.dist, cl.ret[idx])
}

// Clones the repository in other goroutine
func (cl *cloner) clone(repo string) {
	for i, done := range cl.done {
		if done {
			cl.cloneWith(i, repo)
			return
		}
	}
	cl.cloneWith(cl.waitOne(), repo)
}
