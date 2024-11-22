package goconcurrentpool

import (
	"fmt"
	"sync"
)

type job func() (any, error)

type concurrentPool struct {
	size      int
	jobChan   chan job
	closeChan chan struct{} // Канал для отслеживания состояния заверешния работы
	active    bool
	wg        sync.WaitGroup
	mu        sync.Mutex
}

func New(size int) *concurrentPool {
	if size <= 0 {
		size = 1
	}
	return &concurrentPool{
		size:      size,
		jobChan:   make(chan job),
		closeChan: make(chan struct{}),
	}
}

func (p *concurrentPool) Run() {

	p.mu.Lock()
	p.active = true
	p.mu.Unlock()

	for i := 0; i < p.size; i++ {
		go func() {
			for jobFunc := range p.jobChan {
				func() { // локализация обработки паники при выполнении jobFunc
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("concurrent pool: jobFunc recovered from panic: %+v", r)
						}
					}()
					jobFunc()
				}()
				p.wg.Done()
			}
		}()
	}
}

func (p *concurrentPool) RunJob(jobFunc job) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.active {
		// Канал закрыт, не добавляем новые задания
		return ErrPoolNotActive
	}

	// Канал открыт, добавляем задание
	p.wg.Add(1)
	p.jobChan <- jobFunc

	return nil
}

func (p *concurrentPool) WaitAndClose() {
	p.mu.Lock()
	p.active = false
	p.mu.Unlock()

	p.wg.Wait()
	close(p.jobChan)
}
