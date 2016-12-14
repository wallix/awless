package aws

import (
	"fmt"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	AccessService *Access
	InfraService  *Infra
)

func InitSession(region string) (*session.Session, error) {
	session, err := session.NewSession(&awssdk.Config{Region: awssdk.String(region)})
	if err != nil {
		return nil, err
	}
	if _, err = session.Config.Credentials.Get(); err != nil {
		return nil, fmt.Errorf("Your AWS credentials seem undefined: %s", err)
	}

	return session, nil
}

func InitServices(sess *session.Session) {
	AccessService = NewAccess(sess)
	InfraService = NewInfra(sess)
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
