package main

import (
	"fmt"
	"sync"
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

	// Start worker goroutine
	go replyJobs()

	// Check how many times iterated
	sum := 0
	for range replyChan { // Receiver
		sum += 1
	}
	assert.Equal(t, 5, sum)

	// PANIC! send on closed channel!
	// replyChan <- 6
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

func TestBufferedChannel(t *testing.T) {
	// We need to wait for goroutines to finish.
	var wg sync.WaitGroup

	// Buffered channels are useful when you know how many goroutines you have launched,
	// want to limit the number of goroutines you will launch,
	// or want to limit the amount of work that is queued up.
	//
	// When a sender sends a value on a buffered channel, it blocks only if the channel is full.
	replyChan := make(chan int, 2)

	wg.Add(1) // Add 1 goroutine
	go func(ch chan int) {
		defer wg.Done()

		ch <- 1
		ch <- 2

		// DEADLOCK!
		// ch <- 3
	}(replyChan)

	wg.Wait()
	close(replyChan)
}

func TestPanicSituations(t *testing.T) {
	// Deadlock
	//
	// A deadlock happens when all goroutines involved in a concurrent program are blocked,
	// waiting for each other to proceed, resulting in a situation where no progress can be made.
	// Deadlocks commonly occur in Go when there is a mismatch in the sending and receiving operations on channels.
	//
	//

	/**
	// Deadlock 1: Empty channel (No receiver)
	// No other groutine running

	valChan := make(chan int)

	// Do nothing spawned goroutine
	go func() {}()

	fmt.Println("DEADLOCK 1", <-valChan)
	*/

	/**
	// Deadlock 2: Empty channel (No sender)
	// No other groutine running

	valChan2 := make(chan int)
	go func() {
		// Send data to the channel (i)
		valChan2 <- 10
		//
		// No other goroutines that can send to the channel either
		// valChan2 <- 20 // (Send data to the channel (j))
	}()

	_ = <-valChan2 // (OK)
	// PANIC: DEADLOCK 2:
	// Tries to receive another value, but no other groutine running
	_ = <-valChan2
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

func TestChannelBuffered(t *testing.T) {
	fmt.Println("------------------- TestChannelBuffered -------------------")

	// Executions like queue task
	replyChan := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
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

func TestChannelUnbuffered(t *testing.T) {
	fmt.Println("------------------- TestChannelUnbuffered -------------------")

	// Executions like queue task
	replyChan := make(chan int, 3)
	go func() {
		for i := 0; i < 5; i++ {
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
