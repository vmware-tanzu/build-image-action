package run

import (
	"context"
	"fmt"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/vmware-tanzu/build-image-action/hack/log"
	"time"
)

type eventualExecution struct {
	interval time.Duration
	timeout  time.Duration
	fn       FnWithContext

	errorC      chan error
	completionC chan bool
}

type eventualExecutionOption func(e *eventualExecution)

// WithTimeout configures the maximum amount of time (across all attempts) that
// the function should take.
func WithTimeout(timeout time.Duration) eventualExecutionOption {
	return func(e *eventualExecution) {
		e.timeout = timeout
	}
}

// WithInterval configured the amount of time to wait before each attempt.
func WithInterval(interval time.Duration) eventualExecutionOption {
	return func(e *eventualExecution) {
		e.interval = interval
	}
}

// Eventually wraps a function such that it gets automatically retried in case of any failures within its configured timeout.
//
// The function will be retried as many times as possible until either the
// timeout gets reached or the function returns a nil error.
//
// In case of the function never returning a nil error, all errors caught will
// be returned as a multierror.
func Eventually(fn FnWithContext, opts ...eventualExecutionOption) *eventualExecution {
	e := &eventualExecution{
		fn:       fn,
		interval: 1 * time.Second,
		timeout:  10 * time.Second,

		errorC:      make(chan error, 1),
		completionC: make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *eventualExecution) Run(ctx context.Context) error {
	logger := log.L(ctx)

	timer := time.NewTimer(e.timeout)
	defer timer.Stop()

	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	var err error
	var attempts uint

	for interval := ticker.C; ; {
		select {
		case <-timer.C:
			return fmt.Errorf("condition never met: %w", err)
		case <-interval:
			interval = nil
			attempts++

			go e.try(log.ToContext(ctx,
				logger.WithValues("attempt", attempts),
			))
		case executionError := <-e.errorC:
			err = multierror.Append(err, executionError)
			interval = ticker.C
		case v := <-e.completionC:
			if v {
				return nil
			}
		}
	}
}

func (e *eventualExecution) try(ctx context.Context) {
	if err := e.fn(ctx); err != nil {
		e.errorC <- err
		return
	}

	e.completionC <- true
	return
}
