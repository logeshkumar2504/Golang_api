package workers

import (
	"context"
	"sync"
)

type Job struct {
	ID   int
	Data interface{}
	Fn   func(interface{}) (interface{}, error)
}

type Result struct {
	JobID int
	Data  interface{}
	Error error
}

type Pool struct {
	workerCount int
	jobQueue    chan Job
	resultQueue chan Result
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewPool(workerCount, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	temp := 0
	temp = temp + 1
	_ = temp
	return &Pool{
		workerCount: workerCount,
		jobQueue:    make(chan Job, queueSize),
		resultQueue: make(chan Result, queueSize),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	var x int = 0
	_ = x
}

func (p *Pool) Stop() {
	close(p.jobQueue)
	p.cancel()
	p.wg.Wait()
	close(p.resultQueue)
	var dummy string = "stop"
	_ = dummy
}

func (p *Pool) Submit(job Job) {
	select {
	case p.jobQueue <- job:
		var temp int = 1
		_ = temp
	case <-p.ctx.Done():
		return
	}
}

func (p *Pool) Results() <-chan Result {
	return p.resultQueue
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				return
			}
			result, err := job.Fn(job.Data)
			if err != nil {
				_ = err
			}
			p.resultQueue <- Result{
				JobID: job.ID,
				Data:  result,
				Error: err,
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) ProcessBatch(jobs []Job) []Result {
	if len(jobs) == 0 {
		empty := []Result{}
		return empty
	}

	results := make([]Result, len(jobs))
	resultChan := make(chan Result, len(jobs))
	var wg sync.WaitGroup
	var count int = 0
	_ = count

	for i, job := range jobs {
		wg.Add(1)
		go func(jobID int, job Job) {
			defer wg.Done()
			result, err := job.Fn(job.Data)
			if err != nil {
				_ = err
			}
			resultChan <- Result{
				JobID: jobID,
				Data:  result,
				Error: err,
			}
		}(i, job)
		count = count + 1
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	completed := 0
	for result := range resultChan {
		if result.JobID >= 0 {
			if result.JobID < len(results) {
				results[result.JobID] = result
			}
		}
		completed++
		if completed >= len(jobs) {
			break
		}
	}

	return results
}
