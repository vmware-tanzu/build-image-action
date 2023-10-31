// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/stern/stern/stern"
	"k8s.io/client-go/kubernetes"
)

var tails = make(map[string]*stern.Tail)
var tailLock sync.RWMutex

func getTail(targetID string) (*stern.Tail, bool) {
	tailLock.RLock()
	defer tailLock.RUnlock()
	tail, ok := tails[targetID]
	return tail, ok
}

func setTail(targetID string, tail *stern.Tail) {
	tailLock.Lock()
	defer tailLock.Unlock()
	tails[targetID] = tail
}

func clearTail(targetID string) {
	tailLock.Lock()
	defer tailLock.Unlock()
	delete(tails, targetID)
}

// Run starts the main run loop
func Run(ctx context.Context, clientSet *kubernetes.Clientset, config *stern.Config) error {
	if len(config.Namespaces) > 1 {
		return fmt.Errorf("only single namespace supported, got %s", config.Namespaces)
	}

	added := make(chan *stern.Target)
	removed := make(chan *stern.Target)
	errCh := make(chan error)

	defer close(added)
	defer close(removed)
	defer close(errCh)

	a, r, err := stern.Watch(ctx,
		clientSet.CoreV1().Pods(config.Namespaces[0]),
		config.PodQuery,
		config.ExcludePodQuery,
		config.ContainerQuery,
		config.ExcludeContainerQuery,
		config.InitContainers,
		config.EphemeralContainers,
		config.ContainerStates,
		config.LabelSelector,
		config.FieldSelector)
	if err != nil {
		return errors.Wrap(err, "failed to set up watch")
	}

	go func() {
		for {
			select {
			case v, ok := <-a:
				if !ok {
					errCh <- fmt.Errorf("lost watch connection")
					return
				}
				added <- v
			case v, ok := <-r:
				if !ok {
					errCh <- fmt.Errorf("lost watch connection")
					return
				}
				removed <- v
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for p := range added {
			targetID := p.GetID()

			if tail, ok := getTail(targetID); ok {
				if tail.IsActive() {
					continue
				}
				// fmt.Printf("::endgroup::\n")
				tail.Close()
				clearTail(targetID)
			}

			// fmt.Printf("::group::%s\n", p.Container)

			tail := stern.NewTail(clientSet.CoreV1(), p.Node, p.Namespace, p.Pod, p.Container, config.Template, config.Out, config.ErrOut, &stern.TailOptions{
				Timestamps:   config.Timestamps,
				Location:     config.Location,
				SinceSeconds: int64(config.Since.Seconds()),
				Exclude:      config.Exclude,
				Include:      config.Include,
				Namespace:    config.AllNamespaces || len(config.Namespaces) > 1,
				TailLines:    config.TailLines,
				Follow:       config.Follow,
			})
			setTail(targetID, tail)

			go func(tail *stern.Tail) {
				if err := tail.Start(ctx); err != nil {
					fmt.Fprintf(config.ErrOut, "unexpected error: %v\n", err)
				}
			}(tail)
		}
	}()

	go func() {
		for p := range removed {
			targetID := p.GetID()
			if tail, ok := getTail(targetID); ok {
				tail.Close()
				clearTail(targetID)
			}
		}
	}()

	select {
	case e := <-errCh:
		return e
	case <-ctx.Done():
		return nil
	}
}
