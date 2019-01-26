package worker

import (
	"github.com/rs/zerolog/log"
	"sync/atomic"
)

// atomic job counter
var jobCounter uint64

// Job represents the job to be run
type Job interface {
	Name() string
	Execute() error
}

// Pool represents a set of similar jobs
type Pool struct {
	JobQueue chan Job
	quit     chan bool
}

// NewPool creates a new pool for jobs with workers
func NewPool(maxWorkers, queueSize int) *Pool {
	pool := Pool{
		JobQueue: make(chan Job, queueSize),
		quit:     make(chan bool),
	}

	for i := 0; i < maxWorkers; i++ {
		go pool.Process()
	}

	return &pool
}

// Process method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (p *Pool) Process() {
	for {
		select {
		case job, ok := <-p.JobQueue:
			if !ok {
				name := ""
				if job != nil {
					name = job.Name()
				}
				log.Error().Msgf("WORKER:: Stopping dequeuing job for : %s", name)
				continue
			}
			log.Info().Msgf("WORKER:: processing job for : %s", job.Name())

			// we have received a work request, increment the counter
			atomic.AddUint64(&jobCounter, 1)
			if err := job.Execute(); err != nil {
				log.Error().Err(err).Msgf("WORKER:: Error executing job for : %s", job.Name())
			}

			// done with the job, decrement counter
			atomic.AddUint64(&jobCounter, ^uint64(0))
			log.Info().Msgf("WORKER:: job done for : %s", job.Name())

		case <-p.quit:
			// we have received a signal to stop
			return
		}
	}
}

// Queue queues up the job received
func (p *Pool) Queue(job Job) {
	go func() {
		p.JobQueue <- job
	}()
}

// Stop signals the pool to stop listening for work requests
func (p *Pool) Stop() {
	// close channels
	close(p.JobQueue)

	go func() {
		p.quit <- true
	}()
}

// Count give the current number of jobs in queue
func (p *Pool) Count() uint64 {
	return atomic.LoadUint64(&jobCounter)
}
