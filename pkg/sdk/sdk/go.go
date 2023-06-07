package sdk

import (
	"context"
	"fmt"
	"sync"
)

type callback = func() error

func Go(fns []callback) (error, context.CancelFunc) {
	errCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	if len(fns) == 0 {
		close(errCh)
		return nil, cancel
	}

	wg := &sync.WaitGroup{}
	for i, fn := range fns {
		wg.Add(1)
		go func(current int, fn callback, wg *sync.WaitGroup) {
			select {
			case <-ctx.Done():
				fmt.Println("all cancel")
				wg.Done()
				errCh <- nil

				return
			default:
				err := fn()

				if err != nil {
					cancel()
				}

				wg.Done()
				errCh <- err

				return
			}
		}(i, fn, wg)
	}

	wg.Wait()

	return <-errCh, cancel
}
