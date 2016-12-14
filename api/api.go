package api

import (
	"sync"

	"github.com/wallix/awless/config"
)

var (
	AccessService *Access
	InfraService  *Infra
)

func InitServices() {
	AccessService = NewAccess(config.Session)
	InfraService = NewInfra(config.Session)
}

func multiFetch(fns ...func() (interface{}, error)) (<-chan interface{}, <-chan error) {
	resultc := make(chan interface{})
	errc := make(chan error, 1)

	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)
		go func(fetchFn func() (interface{}, error)) {
			defer wg.Done()
			r, err := fetchFn()
			if err != nil {
				errc <- err
				return
			}
			resultc <- r
		}(fn)
	}

	go func() {
		wg.Wait()
		close(resultc)
		close(errc)
	}()

	return resultc, errc
}
