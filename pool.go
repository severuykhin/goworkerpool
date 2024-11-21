package goconcurrentpool

import (
	"fmt"
	"sync"
)

type job func()

type concurrentPool struct {
	size      int
	jobChan   chan job
	closeChan chan struct{} // Канал для отслеживания состояния заверешния работы
	wg        sync.WaitGroup
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
	for i := 0; i < p.size; i++ {
		go func() {
			for jobFunc := range p.jobChan {
				func() { // локализация обработки паники при выполнении jobFunc
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("concurent pool: jobFunc recovered from panic: %+v", r)
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
	select {
	case <-p.closeChan:
		// Канал закрыт, не добавляем новые задания
		return fmt.Errorf("concurrent pool is stopped")
	default:
		// Канал открыт, добавляем задание
		p.wg.Add(1)
		p.jobChan <- jobFunc
	}
	return nil
}

func (p *concurrentPool) WaitAndClose() {
	close(p.closeChan) // Закрываем канал состояния
	p.wg.Wait()
	close(p.jobChan)
}
