package run_test

import (
	"context"
	"github.com/vmware-tanzu/build-image-action/hack/run"
	"sync/atomic"
	"testing"
	"time"
)

func TestEventually_attemptsUntilCompletion(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var runs uint64
	fn := func(ctx context.Context) error {

		// there is some use of goroutines behind the scenes, so, be
		// sure that we're not going to end up racing.
		atomic.AddUint64(&runs, 1)

		// simulating a "slow" function; this allows us to validate our
		// desire of having the function being executed until
		// completion rather then on every interval seeing a new run of
		// the same function.
		time.Sleep(1 * time.Second)
		return nil
	}

	if err := run.Eventually(fn,
		run.WithTimeout(3*time.Second),
		run.WithInterval(50*time.Millisecond),
	).Run(ctx); err != nil {
		t.Fatal(err)
	}

	expectedRuns := uint64(1)
	if expectedRuns != runs {
		t.Fatalf("expected %d runs, had %d", expectedRuns, runs)
	}
}
