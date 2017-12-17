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
	slugs   chan string
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
		slugs:   make(chan string, maxBuffer),
	}

	if c.git == "" {
		c.git = "git"
	}

	return c
}

// Clone clones the repository. Format of 'slug' parameter is 'owner/name'.
func (cl *Cloner) Clone(slug string) {
	cl.slugs <- slug
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
		for slug := range cl.slugs {
			log.Println("Cloning", slug)

			var url string
			if cl.SSH {
				url = fmt.Sprintf("git@github.com:%s.git", slug)
			} else {
				url = fmt.Sprintf("https://github.com/%s.git", slug)
			}

			dir := filepath.FromSlash(fmt.Sprintf("%s/%s", dest, slug))
			cmd := exec.Command(git, "clone", "--depth=1", "--single-branch", url, dir)
			err := cmd.Run()

			if err != nil {
				log.Println("Failed to clone", slug, err)
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
					log.Println("Failed to extract files", slug, extract.String(), err)
					if cl.Err != nil {
						cl.Err <- err
					}
					return
				}
			}

			log.Println("Cloned:", slug)
		}
	}()
}

// Start starts underlying workers and makes ready for running.
// Parameter 'para' indicates how many workers should be used.a
// Max number is '# of CPU - 1' and 0 indicates using the default value.
func (cl *Cloner) Start(para int) {
	auto := runtime.NumCPU() - 1
	if para == 0 || para > auto {
		para = auto
	}
	log.Println("Start to clone with", para, "workers")
	for i := 0; i < para; i++ {
		cl.newWorker()
	}
}

// Shutdown stops all workers and waits until all of current tasks are completed.
func (cl *Cloner) Shutdown() {
	close(cl.slugs)
	cl.wg.Wait()
	if cl.Err != nil {
		close(cl.Err)
	}
}
