package wp

import "sync"

type Task interface {
	Execute()
}

/*
	Mutex to avoid goroutines leak
	Size is number of workers processing tasks
	tasks is a task pool
	kill channel to close a goroutine worker
	wg wrapper to wait until all tasks will be processed
 */
type Pool struct {
	mutex sync.Mutex
	size int
	tasks chan Task
	kill chan struct{}
	wg sync.WaitGroup
}

func NewPool(size int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, 128),
		kill: make(chan struct{}),
	}

	pool.Resize(size)

	return pool
}

func (p *Pool) Resize(n int) {
	p.mutex.Lock()

	defer p.mutex.Unlock()

	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}

	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()

	for {
		select {
			case task, ok := <- p.tasks:
				if ! ok {
					return
				}

				task.Execute()
			case <-p.kill:
				return
		}
	}
}

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Exec(task Task) {
	p.tasks <- task
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
