package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
)

// pullers is an arbitary limit of how many pullers can be running at the same time, so that we don't
// hit the rate limits of source-code hosting provider API.
const pullers = 4

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get current path: %v", err))
	}

	dirs, err := getAllDirs(currentPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get list of directories: %v", err))
	}

	if len(dirs) == 0 {
		log.Fatal("no directories found")
	}

	work := make(chan string)
	wg := sync.WaitGroup{}
	for i := 0; i < pullers; i++ {
		wg.Add(1)
		go func(ctx context.Context, num int) {
			defer wg.Done()

			log.Printf("puller: %v, starting\n", num)
			for {
				select {
				case <-ctx.Done():
					{
						log.Printf("puller: %v, shutting down\n", num)
					}
				case dir, ok := <-work:
					if !ok {
						log.Printf("puller: %v, finished\n", num)
						return
					}

					err := fetchAndPull(dir)
					if err != nil {
						log.Printf("puller: %v, failed to fetch and pull repo: %v, error: %v", num, dir, err)
					} else {
						log.Printf("puller: %v, updated repo: %v", num, dir)
					}

				}
			}
		}(ctx, i)
	}

	for _, dir := range dirs {
		work <- dir.Name()
	}
	close(work)

	waitUntilDoneOrShutdown(cancel, &wg)
}

func waitUntilDoneOrShutdown(cancel context.CancelFunc, s *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	s.Wait()
}

func fetchAndPull(dir string) error {
	comms := [][]string{
		{"git", "checkout", "main"},
		{"git", "fetch"},
		{"git", "pull", "--ff-only"},
	}

	for _, com := range comms {
		cmd := exec.Command(com[0], com[1:]...)
		cmd.Dir = dir
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run command: %v, error: %v", com, err)
		}
	}

	return nil
}

// getAllDirs returns all directories in the given path.
func getAllDirs(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %v", err)
	}

	onlyDirs(&files)

	return files, nil
}

// onlyDirs returns only directories from the given list of files.
func onlyDirs(files *[]os.DirEntry) {
	j := 0
	for i := 0; i < len(*files); i++ {
		if (*files)[i].IsDir() {
			(*files)[j] = (*files)[i]
			j++
		}
	}

	*files = (*files)[:j]
}
