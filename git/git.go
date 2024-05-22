package git

import (
	"fmt"
	"sync"
	"time"
)

type Branch struct {
	Name     string
	IsMaster bool
}

type Merge struct {
	To       *Branch
	From     *Branch
	Approved bool
}

type PullRequest struct {
	Actor  string
	Branch Branch
}

type Work struct {
	merge       chan *Merge
	pullRequest chan *PullRequest
}

type Git struct {
	work Work
}

func New() *Git {
	return &Git{
		work: Work{
			merge:       make(chan *Merge),
			pullRequest: make(chan *PullRequest),
		},
	}
}

func (g *Git) Run() {
	g.submitWork(&g.work)
}

func (g *Git) submitWork(w *Work) {
	var wg sync.WaitGroup

	wg.Add(2)
	go g.mergeByOwner(&wg, w.merge, w.pullRequest)
	go g.pullRequest(&wg, w.merge, w.pullRequest)

	g.Shutdown(&wg)
}

func (g *Git) mergeByOwner(wg *sync.WaitGroup, m chan<- *Merge, p <-chan *PullRequest) {
	defer wg.Done()

	// Check PR
	time.Sleep(2 * time.Second) // time processed
	pr := <-p
	if pr.Actor != "" {
		fmt.Println("PR from: ", pr.Actor, " processed")

		// Only owner can merge to master
		time.Sleep(2 * time.Second) // time processed
		m <- &Merge{
			To: &Branch{
				Name:     "Master",
				IsMaster: true,
			},
			From:     &pr.Branch,
			Approved: true,
		}

		fmt.Println("PR merged from: ", pr.Branch.Name, "to: Master ")
	}
}

func (g *Git) pullRequest(wg *sync.WaitGroup, m <-chan *Merge, p chan<- *PullRequest) {
	defer wg.Done()

	// Create pull request
	//
	// Any PR will requested to owner
	fmt.Println("PR requested")
	p <- &PullRequest{
		Actor: "Contributor",
		Branch: Branch{
			Name: "Feature",
		},
	}

	// Check merged request has been merged
	merged := <-m
	if merged.Approved {
		fmt.Println("PR succesfully merged")
	}
}

func (g *Git) Shutdown(wg *sync.WaitGroup) {
	wg.Wait()
	close(g.work.merge)
	close(g.work.pullRequest)
}
