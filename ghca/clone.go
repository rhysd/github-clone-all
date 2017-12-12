package ghca

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
)

const maxConcurrency = 4
const maxBuffer = 1000

type Cloner struct {
	git     string
	dest    string
	extract *regexp.Regexp
	repos   chan string
	Err     chan error
	wg      sync.WaitGroup
	SSH     bool
}

func NewCloner(dest string, extract *regexp.Regexp) *Cloner {
	c := &Cloner{
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

func (cl *Cloner) Clone(repo string) {
	cl.repos <- repo
}

func (cl *Cloner) newWorker() {
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
			if cl.SSH {
				url = fmt.Sprintf("git@github.com:%s.git", repo)
			} else {
				url = fmt.Sprintf("https://github.com/%s.git", repo)
			}

			dir := filepath.FromSlash(fmt.Sprintf("%s/%s", dest, repo))
			cmd := exec.Command(git, "clone", "--depth=1", "--single-branch", url, dir)
			err := cmd.Run()

			if err != nil {
				log.Println("Failed to clone", repo, err)
				if cl.Err != nil {
					cl.Err <- err
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
					if (info.Mode()&os.ModeSymlink != 0) || !extract.MatchString(path) {
						return os.Remove(path)
					}
					return nil
				}); err != nil {
					log.Println("Failed to extract files", repo, extract.String(), err)
					if cl.Err != nil {
						cl.Err <- err
					}
					return
				}
			}

			log.Println("Cloned:", repo)
		}
	}()
}

func (cl *Cloner) Start() {
	para := runtime.NumCPU() - 1
	log.Println("Start to clone with", para, "workers")
	for i := 0; i < para; i++ {
		cl.newWorker()
	}
}

func (cl *Cloner) Shutdown() {
	close(cl.repos)
	cl.wg.Wait()
	if cl.Err != nil {
		close(cl.Err)
	}
}
