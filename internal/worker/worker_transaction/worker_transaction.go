package workertransaction

import (
	"errors"
	"sync"
)

var (
	ErrNoJob     = errors.New("no job available")
	counterError = 0
)

type workerCleaner struct {
	wg *sync.WaitGroup
	WorkerCleanerConfig
	// dw dwc.Repository
}

type WorkerCleanerConfig struct {
	NumWorker, WorkerDelay, RetryDelay, RetryAttempt uint
	Schedule                                         string
}
