package main

import (
	"fmt"
	"sync"
)

type Task func()

type ThreadPool struct {
	queue      chan Task
	numWorkers int
	wg         sync.WaitGroup
}

func NewThreadPool(tasksLimit int, numWorkers int) *ThreadPool {
	tp := &ThreadPool{
		queue:      make(chan Task, tasksLimit),
		numWorkers: numWorkers,
	}
	go tp.run()
	return tp
}

func (tp *ThreadPool) run() {
	for id := 1; id <= tp.numWorkers; id++ {
		tp.wg.Add(1)
		go tp.worker(id)
	}
}

func (tp *ThreadPool) worker(id int) {
	defer fmt.Printf("Worker %d Shutdown\n", id)
	defer tp.wg.Done()
	for task := range tp.queue {
		fmt.Printf("Worker %d\n", id)
		task()
	}
}

func (tp *ThreadPool) submit(task Task) {
	if cap(tp.queue) <= len(tp.queue) {
		fmt.Println("Queue is overloaded. Try again")
		return
	}

	tp.queue <- task
	fmt.Printf("Pushed new task to Queue, task count: %d\n", len(tp.queue))
}

func (tp *ThreadPool) wait() {
	tp.wg.Wait()
}

func (tp *ThreadPool) shutdown() {
	close(tp.queue)
}

func main() {
	threadPool := NewThreadPool(10, 2)
	for i := range 20 {
		threadPool.submit(func() { fmt.Printf("TASK %d\n", i) })
	}
	threadPool.shutdown()
	threadPool.wait()
}
