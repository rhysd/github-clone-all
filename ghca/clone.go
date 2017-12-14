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

// Cloner is a git-clone worker to clone given repositories with workers in parallel.
type Cloner struct {
	git     string
	dest    string
	extract *regexp.Regexp
	repos   chan string
	// Err is a receiver of errors which occurs while cloning repositories
	Err chan error
	wg  sync.WaitGroup
	// SSH is a flag to use SSH for git-clone. By default, it's false and HTTPS is used.
	SSH bool
}

// NewCloner creates a new cloner instance. 'extract' parameter can be nil.
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

// Clone clones the repository. Format of 'repo' parameter is 'owner/name'.
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

// Start starts underlying workers and makes ready for running.
func (cl *Cloner) Start() {
	para := runtime.NumCPU() - 1
	log.Println("Start to clone with", para, "workers")
	for i := 0; i < para; i++ {
		cl.newWorker()
	}
}

// Shutdown stops all workers and waits until all of current tasks are completed.
func (cl *Cloner) Shutdown() {
	close(cl.repos)
	cl.wg.Wait()
	if cl.Err != nil {
		close(cl.Err)
	}
}
