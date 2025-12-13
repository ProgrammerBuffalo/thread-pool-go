package main

import (
	"fmt"
	"sync"
)

type Task func()

type ThreadPool struct {
	queue chan Task
	sem   chan struct{}
	wg    sync.WaitGroup
}

func NewThreadPool(tasksLimit int, threadLimit int) *ThreadPool {
	tp := &ThreadPool{
		queue: make(chan Task, tasksLimit),
		sem:   make(chan struct{}, threadLimit),
	}
	go tp.run()
	return tp
}

func (tp *ThreadPool) run() {
	for id := 0; ; id++ {
		select {
		case task := <-tp.queue:
			tp.sem <- struct{}{}
			go func(id int) {
				fmt.Printf("Thread %d acquired\n", id)
				defer func() {
					<-tp.sem
					fmt.Printf("Thread %d released\n", id)
					tp.wg.Done()
				}()
				task()
			}(id)
		}
	}
}

func (tp *ThreadPool) submit(task Task) {
	if cap(tp.queue) <= len(tp.queue) {
		fmt.Println("Queue is overloaded. Try again")
		return
	}

	tp.wg.Add(1)
	tp.queue <- task
	fmt.Printf("Pushed new task to Queue, task count: %d\n", len(tp.queue))
}

func (tp *ThreadPool) wait() {
	tp.wg.Wait()
}

func main() {
	threadPool := NewThreadPool(10, 2)
	for i := 0; i < 20; i++ {
		threadPool.submit(func() { fmt.Printf("====TASK %d=====", i) })
	}
	threadPool.wait()
}
