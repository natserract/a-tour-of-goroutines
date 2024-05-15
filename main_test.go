package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicDoneChan(t *testing.T) {
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
	t.Parallel()

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
	t.Parallel()

	// Closing a channel indicates that no more values will be sent on it.
	// A closed channel states that we can't send data to it, but we can still read data from it.
	replyChan := make(chan int)

	// Sender
	replyJobs := func() {
		// Send items to be processed
		// Wait all goroutines to finish
		for i := 0; i < 5; i++ {
			replyChan <- i
		}

		// Only the sender should close the channel
		// Sending data to a closed channel will panic. So to ensure that the receiver doesn’t prematurely
		// close the channel while the sender is still sending data to it, the sender should close the channel.
		close(replyChan)
	}

	// Start goroutines
	go replyJobs()

	// Check how many times iterated
	sum := 0
	for range replyChan { // Receiver
		sum += 1
	}
	assert.Equal(t, 5, sum)

	// PANIC!
	// replyChan <- 6
}

func TestCloseWithDeferChan(t *testing.T) {
	t.Parallel()

	replyChan := make(chan int)
	replySender := func() {
		// A defer statement defers the execution of a function until the surrounding function returns.
		// No matter how many return in your goroutine, close function is guaranteed to be executed right after the function (replyJobs) returns.
		defer close(replyChan)

		for i := 0; i < 5; i++ {
			if i == 4 {
				// If an error or panic occurs before reaching the
				// end of the function, the channel will still be closed
				// due to the deferred close operation.
				panic("Something went wrong!") // Simulate a panic occurring during sendData
			} else {
				replyChan <- i
			}
		}
	}

	replyReceiver := func() {
		for val := range replyChan {
			fmt.Println("Fetched val: ", val)
		}
	}

	// Start 1 goroutine
	go replySender()
	replyReceiver()
	fmt.Println("Next line will executed")
}
