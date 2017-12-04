package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type cloner struct {
	recv chan string
	git  string
	dist string
}

func newCloner(dist string) *cloner {
	git := os.Getenv("GIT_EXECUTABLE_PATH")
	if git == "" {
		git = "git"
	}
	return &cloner{make(chan string), git, dist}
}

func (cl *cloner) clone(repo string) error {
	url := fmt.Sprintf("https://github.com/%s.git", repo)
	dist := fmt.Sprintf("%s/%s", cl.dist, repo)
	cmd := exec.Command(cl.git, "clone", "--depth=1", "--single-branch", url, dist)
	return cmd.Run()
}

func (cl *cloner) start() {
	for {
		select {
		case repo := <-cl.recv:
			if repo == "" {
				// When channel is closed
				return
			}
			log.Println("Cloning", repo)
			if err := cl.clone(repo); err != nil {
				log.Println("Failed to clone", repo)
			} else {
				log.Println("Cloned", repo)
			}
		}
	}
}
