# Go-routines walkthrough

**Table of contents:**

1. Intro to go channels
2. Channel communication strategy
3. Error handling
4. Shared state, and appending external state
5. Waitgroup & Mutex (Synchronization Primitives)
6. Batching

## 1. Intro to Go Channels

Pengenalan basic tentang go channels

Channels are concurrent-safe queues that are used to safely pass messages between Go’s lightweight processes (goroutines). Together, these primitives are some of the most popularly touted features of the Go programming language. The message-passing style they encourage permits the programmer to safely coordinate multiple concurrent tasks with easy-to-reason-about semantics and control flow that often trumps the use of callbacks or shared memory.

Concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of computation. Parallelism is about executing many things at once, it’s focus is execution. While concurrency is about dealing with many things at once, it’s focus is structure

Channnels:
To communicate and connect between these goroutines by channels

Only the sender should close the channel
Sending data to a closed channel will panic. So to ensure that the receiver doesn’t prematurely close the channel while the sender is still sending data to it, the sender should close the channel.

## 2. Channel communication strategy

Unidirectional Channels,
Bidirectional Channels,
Buffered Channels,
Select Statement,
Closing Channels,
Range over channels
