package http

import (
	execution2 "emulator/pkg/execution"
	_var "emulator/var"
	"sync"
	"time"
)

func CloseExecutioners() {
	wg := sync.WaitGroup{}
	for _, e := range []string{_var.PROJECT_EXECUTION} {
		wg.Add(1)

		go func(name string, wg *sync.WaitGroup) {
			if !execution2.Service(name).Closed() {
				execution2.Service(name).Close()
			}
			wg.Done()
		}(e, &wg)
	}
	wg.Wait()

	time.Sleep(5 * time.Second)
	execution2.FinalCleanup(true)
}
