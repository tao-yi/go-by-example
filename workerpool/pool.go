package workerpool

import "fmt"

type Result struct {
	Err  error
	Data interface{}
}

type Task func() *Result

type Pool struct {
	size        int
	tasks       []Task
	resultQueue chan *Result
	taskQueue   chan Task
}

func New(size int) *Pool {
	p := &Pool{
		size:        size,
		tasks:       make([]Task, 0, size),
		taskQueue:   make(chan Task, size),
		resultQueue: make(chan *Result, size),
	}

	for i := 0; i < size; i++ {
		fmt.Printf("start worker %d\n", i)
		go p.startWorker()
	}
	return p
}

func (p *Pool) startWorker() {
	for task := range p.taskQueue {
		p.resultQueue <- task()
	}
}

func (p *Pool) dispatchTasks() {
	for i, task := range p.tasks {
		fmt.Printf("dequeue task %d\n", i)
		p.taskQueue <- task
	}
	close(p.taskQueue)
	// clear tasks
	p.tasks = make([]Task, 0, p.size)
}

func (p *Pool) Submit(task Task) {
	p.tasks = append(p.tasks, task)
}

func (p *Pool) Await() []*Result {
	taskCount := len(p.tasks)
	results := make([]*Result, 0, taskCount)
	if taskCount == 0 {
		return results
	}
	go p.dispatchTasks()
	for i := 0; i < taskCount; i++ {
		results = append(results, <-p.resultQueue)
	}
	return results
}
