package scheduler

import (
	"context"
	"sync"
	"time"
)

// Scheduler  runs delayed tasks
// Attention deferred tasks will start when the task waiting time exceeds the default Scheduler
// The default Scheduler timeout is 10 seconds
// Time default Scheduler waiting time adds in app config initialization
// mutexes are needed for the safe operation of the goroutine with the mop
// because tests runs some goroutines of Scheduler and test getting panics
type Scheduler struct {
	mx       sync.RWMutex
	taskChan chan SchedulerTask
	tasks    map[int32]SchedulerTask
}

func NewScheduler() *Scheduler {
	tasks := make(map[int32]SchedulerTask, 3072)
	taskChan := make(chan SchedulerTask, 64)
	return &Scheduler{taskChan: taskChan, tasks: tasks, mx: sync.RWMutex{}}
}

type SchedulerTask struct {
	waitTime  time.Time
	task      func()
	key       int32
	isRemoved bool
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case task := <-s.taskChan:
			s.add(task)
		case <-ticker.C:
			s.runDelayedTask()
		}
	}
}

// Write is adding your task in Scheduler map
// the key is needed to avoid duplication of tasks and tracking them
func (s *Scheduler) Write(key int32, task func(), duration time.Duration) {
	newSchedulerTask := SchedulerTask{
		key:      key,
		task:     task,
		waitTime: time.Now().Add(time.Second * duration),
	}
	s.taskChan <- newSchedulerTask
}

// Remove is removing  task from Scheduler map
// the key identifies the pending task
func (s *Scheduler) Remove(key int32) {
	newSchedulerTask := SchedulerTask{
		key:       key,
		isRemoved: true,
	}
	s.taskChan <- newSchedulerTask
}

func (s *Scheduler) Count() int32 {
	s.mx.Lock()
	defer s.mx.Unlock()
	return int32(len(s.tasks))
}

func (s *Scheduler) add(task SchedulerTask) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if task.isRemoved {
		delete(s.tasks, task.key)
		return
	}

	s.tasks[task.key] = task
}

func (s *Scheduler) runDelayedTask() {
	s.mx.Lock()
	defer s.mx.Unlock()
	if len(s.tasks) == 0 {
		return
	}
	for key, v := range s.tasks {
		timeNow := time.Now()
		task := v.task
		if timeNow.After(v.waitTime) {
			task()
			delete(s.tasks, key)
		}
	}
}
