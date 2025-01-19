package utils

import (
	"fmt"
	"time"
)

type RetryFunc func() error

func Retry(retryFunc RetryFunc, maxIterations int, delay time.Duration) error {
	var err error
	for i := 0; i < maxIterations; i++ {
		err = retryFunc()
		if err == nil {
			return nil
		}
		fmt.Printf("Attempt %d/%d failed: %v. Retrying in %v seconds...\n", i+1, maxIterations, err, delay.Seconds())
		time.Sleep(delay)
	}
	return err
}
