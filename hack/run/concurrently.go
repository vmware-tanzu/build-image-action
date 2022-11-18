package run

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type (
	FnWithContext func(ctx context.Context) error
	Fn            func() error
)

// ConcurrentlyWithContext runs a set of functions `fns` concurrently with no
// max concurrency control passing down to each one of the functions the same
// context as the one passed to it.
//
// The first function to return a non-nil error will cancel the group
// execution.
func ConcurrentlyWithContext(ctx context.Context, fns ...FnWithContext) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, f := range fns {
		f := f
		g.Go(func() error {
			return f(ctx)
		})
	}

	return g.Wait()
}

// Concurrently concurrently runs the set of functions passed to it.
//
// The first function to return a non-nil error will cancel the group
// execution.
func Concurrently(fns ...Fn) error {
	var g errgroup.Group

	for _, f := range fns {
		f := f
		g.Go(func() error {
			return f()
		})
	}

	return g.Wait()
}

func Serially(fns ...Fn) error {
	for _, f := range fns {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func SeriallyWithContext(ctx context.Context, fns ...FnWithContext) error {
	for _, f := range fns {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}
