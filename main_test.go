package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestBasicDoneChan(t *testing.T) {
	fmt.Println("------------------- TestBasicDoneChan -------------------")

	// Define the unbuffered channel
	doneCh := make(chan bool)

	// Spawn a goroutine
	go func() {
		time.Sleep(1 * time.Second)
		doneCh <- true
	}()

	// Receive the doneCh channel value
	doneResp := <-doneCh

	// The variable declaration doneResp will block until a message is written to it from the channel
	assert.Equal(t, true, doneResp)
}

func TestBasicPlusOneChan(t *testing.T) {
	fmt.Println("------------------- TestBasicPlusOneChan -------------------")

	// Spawn a goroutine and pass the argument
	plusOneCh := make(chan int)

	// Worker/processor to increase start + 1
	plusOneWorker := func(start int) {
		val := start + 1
		plusOneCh <- val
	}

	go plusOneWorker(1)

	plusOneResp := <-plusOneCh
	assert.Equal(t, 2, plusOneResp)
}

func TestCloseChan(t *testing.T) {
	fmt.Println("------------------- TestCloseChan -------------------")

	// Closing a channel indicates that no more values will be sent on it.
	// A closed channel states that we can't send data to it, but we can still read data from it.
	jobsChan := make(chan int)

	// Sender
	createJobs := func() {
		// Send items to be processed
		// Wait all goroutines to finish
		for i := 0; i < 5; i++ {
			jobsChan <- i
		}

		// Only the sender should close the channel
		// Sending data to a closed channel will panic. So to ensure that the receiver doesn’t prematurely
		// close the channel while the sender is still sending data to it, the sender should close the channel.
		close(jobsChan)
	}

	// Start worker goroutine
	go createJobs()

	// Check how many times iterated
	sum := 0
	for range jobsChan { // Receiver
		sum += 1
	}
	assert.Equal(t, 5, sum)

	// PANIC! send on closed channel!
	// jobsChan <- 6
}

func TestCloseWithDeferChan(t *testing.T) {
	fmt.Println("------------------- TestCloseWithDeferChan -------------------")

	replyChan := make(chan int)
	replySender := func() {
		// A defer statement defers the execution of a function until the surrounding function returns.
		// No matter how many return in your goroutine, close function is guaranteed to be executed right after the function (replyJobs) returns.
		defer close(replyChan)

		for i := 0; i < 5; i++ {
			replyChan <- i
		}
	}

	// Start worker goroutine
	go replySender()

	// Receiver
	for val := range replyChan {
		fmt.Println("Fetched val: ", val)
	}

	// If an error or panic occurs before reaching the
	// end of the function, the channel will still be closed
	// due to the deferred close operation.
	fmt.Println("Next line will executed")
}

func TestPanicSituations(t *testing.T) {
	fmt.Println("------------------- TestPanicSituations -------------------")

	/**
	Deadlock

	A deadlock happens when all goroutines involved in a concurrent program are blocked,
	waiting for each other to proceed, resulting in a situation where no progress can be made.

	Deadlocks commonly occur in Go when there is a mismatch in the sending and receiving operations on channels.

	--------------------------------------------------
	DEADLOCK 1: Empty channel (No receiver)
	No other groutine running

	```go
	valChan := make(chan int)

	// Do nothing spawned goroutine
	go func() {}()

	fmt.Println("DEADLOCK 1", <-valChan)
	```

	--------------------------------------------------
	DEADLOCK 2: Empty channel (No sender)
	No other groutine running

	```go
	valChan2 := make(chan int)
	go func() {
		// Send data to the channel (i)
		valChan2 <- 10
		//
		// No other goroutines that can send to the channel either
		// valChan2 <- 20 // (Send data to the channel (j))
	}()

	_ = <-valChan2 // (OK)

	// PANIC!:
	// Tries to receive another value, but no other groutine running
	_ = <-valChan2
	```
	*/
}

func TestChannelDirection(t *testing.T) {
	fmt.Println("------------------- TestChannelDirection -------------------")

	// By default a channel is bidirectional but you can create a unidirectional channel
	//
	// Bidirectional: var bidirectionalChan chan string // can read from, write to and close()
	// Unidirectional (only be used for sending or receiving):
	// 1. var receiveOnlyChan <-chan string // can read from, but cannot write to or close()
	// 2. var sendOnlyChan chan<- string    // cannot read from, but can write to and close()
	//
	// Unidirectional (sending)
	undirectionSendChan := make(chan int)
	sendOnlyWorker := func(c chan<- int) {
		c <- 10
	}
	go sendOnlyWorker(undirectionSendChan)
	assert.Equal(t, 10, <-undirectionSendChan)
	//
	//
	// Unidirectional (receiving)
	undirectionReceiveChan := make(chan int)
	receiveOnlyWorker := func(ch <-chan int) {
		for num := range ch {
			fmt.Println("Received:", num)
		}
	}
	go receiveOnlyWorker(undirectionReceiveChan)

	// Send data to the channel
	for i := 0; i < 3; i++ {
		undirectionReceiveChan <- i
	}
}

func TestChannelUnbuffered(t *testing.T) {
	fmt.Println("------------------- TestChannelUnbuffered -------------------")

	// Executions like queue task
	replyChan := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			// Send data to the channel will blocking until receivers ready
			replyChan <- i
			fmt.Println("Placed: ", i)
		}

		close(replyChan)
	}()

	for n := range replyChan {
		fmt.Println("Preparing ", n)
		time.Sleep(2 * time.Second) // Time processed
		fmt.Println("Served ", n)
		fmt.Println("")
	}
}

func TestChannelBuffered(t *testing.T) {
	fmt.Println("------------------- TestChannelBuffered -------------------")

	replyChan := make(chan int, 3)
	go func() {
		for i := 0; i < 5; i++ {
			// Send data to the channel will blocking until (buffer size full or empty) and receivers ready
			replyChan <- i
			fmt.Println("Placed: ", i)
		}

		close(replyChan)
	}()

	for n := range replyChan {
		fmt.Println("Preparing ", n)
		time.Sleep(2 * time.Second)
		fmt.Println("Served ", n)
		fmt.Println("")
	}
}

func TestBufferedChannelDeadlock(t *testing.T) {
	fmt.Println("------------------- TestBufferedChannelDeadlock -------------------")

	// We need to wait for goroutines to finish.
	var wg sync.WaitGroup

	// Buffered channels are useful when you know how many goroutines you have launched,
	// want to limit the number of goroutines you will launch,
	// or want to limit the amount of work that is queued up.
	//
	// When a sender sends a value on a buffered channel, it blocks only if the channel is full.
	replyChan := make(chan int, 2)

	wg.Add(1) // Start worker goroutine
	go func(ch chan int) {
		defer wg.Done()

		ch <- 1
		ch <- 2

		// DEADLOCK!
		// ch <- 3
	}(replyChan)

	wg.Wait()
	close(replyChan)

	fmt.Println(<-replyChan) // Process the results
	fmt.Println(<-replyChan) // Process the results
	// fmt.Println(<-replyChan) // Process the results
}

func TestErrorHandling(t *testing.T) {
	fmt.Println("------------------- TestErrorHandling -------------------")

	var eg errgroup.Group
	jobsChan := make(chan int)

	// Create new goroutine
	eg.Go(func() error {
		<-jobsChan // receiver
		return fmt.Errorf("Go to error")
	})

	jobsChan <- 10 // Sender
	close(jobsChan)

	// Waiting all goroutines done
	// If err will returned
	err := eg.Wait()
	if assert.Error(t, err, "") {
		expectedErr := errors.New("Go to error")
		assert.Equal(t, expectedErr, err)
	}
}

func TestRaceConditions(t *testing.T) {
	fmt.Println("------------------- TestRaceConditions -------------------")

	var wg sync.WaitGroup
	var jobsChan = make(chan int)

	var total = 0
	send := func(ch chan int) {
		defer wg.Done() // called by each goroutine when it finishes its work, decrementing the internal counter of the wait group.

		for i := 0; i < 10; i++ {
			ch <- i
			total++
		}
	}

	receive := func() {
		for range jobsChan {
		}
	}

	// Spawn 3 goroutines
	gs := 3
	for i := 0; i < gs; i++ {
		wg.Add(1) // used to add the number of goroutines that need to be waited upon.
		go send(jobsChan)
	}

	go func() {
		wg.Wait() // used to block the execution of the goroutine until all the goroutines have called Done()
		close(jobsChan)
	}()

	receive()
	fmt.Println("total", total) // Output: Undetermined (Race condition), should be 30
}

func TestSyncLocking(t *testing.T) {
	fmt.Println("------------------- TestSyncLocking -------------------")

	// Waitgroup
	// WaitGroup to wait for multiple goroutines to finish their execution.
	// It ensures that the main program doesn’t exit before all goroutines have completed.
	//
	// Mutexes (explicit locking)
	// Mutexes are used to prevent multiple goroutines from accessing shared data simultaneously,
	// while channels enable safe communication between multiple goroutines
	//
	// The Mutex (mutual exclusion lock) is an invaluable resource
	// when synchronizing state across multiple goroutines,
	//
	// Problem: multiple Goroutines need to access a shared piece of state
	var wg sync.WaitGroup
	var mu sync.Mutex
	var jobsChan = make(chan int)

	var total = 0

	// to protect shared resources from concurrent access, preventing race conditions.
	send := func(ch chan int) {
		defer mu.Unlock()
		defer wg.Done() // called by each goroutine when it finishes its work, decrementing the internal counter of the wait group.

		mu.Lock()
		for i := 0; i < 10; i++ {
			ch <- i
			total++
		}
	}

	receive := func() {
		for range jobsChan {
		}
	}

	// Spawn 3 goroutines
	gs := 3
	for i := 0; i < gs; i++ {
		wg.Add(1) // used to add the number of goroutines that need to be waited upon.
		go send(jobsChan)
	}

	go func() {
		wg.Wait() // used to block the execution of the goroutine until all the goroutines have called Done()
		close(jobsChan)
	}()

	receive()
	fmt.Println("total", total) // Output: Always 30
	assert.Equal(t, 30, total)
}
