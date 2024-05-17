# Go-routines walkthrough

**Table of contents:**

1. Intro to go channels
2. Channel communication strategy
3. Deadlock situations in the Golang Channel
4. Error handling
5. Shared state, and appending external state
6. Waitgroup & Mutex (Synchronization Primitives)
7. Batching

## 1. Intro to Go Channels

Pengenalan basic tentang go channels

Channels are concurrent-safe queues that are used to safely pass messages between Go’s lightweight processes (goroutines). Together, these primitives are some of the most popularly touted features of the Go programming language. The message-passing style they encourage permits the programmer to safely coordinate multiple concurrent tasks with easy-to-reason-about semantics and control flow that often trumps the use of callbacks or shared memory.

Concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of computation. Parallelism is about executing many things at once, it’s focus is execution. While concurrency is about dealing with many things at once, it’s focus is structure

Channnels:
To communicate and connect between these goroutines by channels.

Golang channel makes goroutines can communicate each other. Through channel, goroutines send or receive messages (values).

Only the sender should close the channel
Sending data to a closed channel will panic. So to ensure that the receiver doesn’t prematurely close the channel while the sender is still sending data to it, the sender should close the channel.

Unbuffered Go Channels Lifecycle (unidirectional):

1. Allocates the channel

```go
make(chan T)
```

2. Wait to senders ready

```go
go func(ch T) {
	ch <- 1
}(ch)
```

3. Wait to receivers ready

It will block until a sender is ready to send the value.

```go
fmt.Println(<-ch)
```

4. Closing

Closing a channel indicates that no more values will be sent on it.

```go
close(ch)
```

Blocking:
There is no receiver for stopChannel<- sender and there is no sender for <-c receiver. In other words: the other side is not ready for each of these channels!

Unbuffered Channels: Sending or receiving on an unbuffered channel blocks the goroutine until the other side is ready. This is useful for ensuring that two goroutines synchronize on some event.

Buffered Channels: These have a buffer of a fixed size and will only block when the buffer is full (sending) or empty (receiving).

```sh
Unbuffered channel: synchronous
Buffered channel: synchronous or asynchronous
```

Every goroutines is independent (concurrent):

```go
channel2 := make(chan int)

	// Sender
	go func() {
		time.Sleep(time.Second)
		channel2 <- 1
	}()
	go func() {
		time.Sleep(time.Second)
		channel2 <- 2
	}()
	go func() {
		time.Sleep(time.Second)
		channel2 <- 3
	}()

	// Receiver
	fmt.Println("Directional Chan", <-channel2)
	fmt.Println("Directional Chan", <-channel2)
	fmt.Println("Directional Chan", <-channel2)

	// Output is dynamically, sometimes different (without wait time.Sleep)
	// Output: 1, 2, 3
	// Output: 2, 1, 3
```

## 2. Channel communication strategy

### Unidirectional Channels,

By default, channels in Go are unidirectional, meaning they can either be used for sending values (<-chan) or receiving values (chan<-). Unidirectional channels enforce the restriction that a channel can only be used for sending or receiving, which can help in making the intent of code clearer.

### Bidirectional Channels,

A channel can also be created as bidirectional (chan). This allows both sending and receiving operations on the channel. Bidirectional channels are useful when you want to use a single channel for both sending and receiving in different parts of your code.

Buffered Channels,
Select Statement,
Closing Channels,
Range over channels
