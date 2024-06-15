package data

import (
	"geth-cody/io"
	"sync"
	"sync/atomic"
)

func RunUntilClosed[T any](bufferSize int, f func(T)) (chan T, func()) {
	c := make(chan T, bufferSize)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range c {
			f(v)
		}
	}()

	return c, func() {
		close(c)
		wg.Wait()
	}
}

type Chan[T any] struct {
	channel chan T
	size    atomic.Int64
}

type Opt[T any] func(*Chan[T])

func WithChannel[T any](channel chan T, size int64) Opt[T] {
	return func(c *Chan[T]) {
		c.channel = channel
		c.size.Store(size)
	}
}

func NewChan[T any](bufferSize int, opts ...Opt[T]) *Chan[T] {
	c := &Chan[T]{
		channel: make(chan T, bufferSize),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Chan[T]) Send(value T) {
	c.size.Add(1)
	c.channel <- value
}

func (c *Chan[T]) Close() {
	close(c.channel)
}

func (c *Chan[T]) Size() int64 {
	return c.size.Load()
}

// AsyncForEach runs f for each value in the channel.
func (c *Chan[T]) ForEach(f func(T) io.Error) io.Error {
	var errs io.ZapErrors
	for v := range c.channel {
		c.size.Add(-1)
		if err := f(v); err != nil {
			errs = append(errs, err)
		}
	}

	return &errs
}

// AsyncForEach runs runs f for each value in the channel over util.Env.ThreadCount different workers.
func (c *Chan[T]) AsyncForEach(bufferSize, threadCount int, f func(T) io.Error) io.Error {
	var errs io.ZapErrors
	errChan, closeErrChan := RunUntilClosed(bufferSize, func(err io.Error) { errs = append(errs, err) })

	var wg sync.WaitGroup
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.ForEach(f); err != nil {
				errChan <- err
			}
		}()
	}

	wg.Done()
	closeErrChan()

	return &errs
}
