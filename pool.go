package goworkerpool

import "sync"

type job func()

type workerPool struct {
	pool chan job
	wg   sync.WaitGroup
}

func NewWorkerPool(size int) *workerPool {
	p := workerPool{
		pool: make(chan job),
	}

	for i := 0; i < size; i++ {
		go func() {
			for wf := range p.pool {
				wf()
				p.wg.Done()
			}
		}()
	}

	return &p
}

func (p *workerPool) RunJob(jobFunc job) {
	p.wg.Add(1)
	p.pool <- jobFunc
}

func (p *workerPool) StopAndWait() {
	p.wg.Wait()
	close(p.pool)
}
