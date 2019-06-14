package schduler

import "sync"

type Schduler struct {
	taskChan chan Task
	taskList []Task
	wg *sync.WaitGroup
}

func NewSchduler() *Schduler  {
	s :=  &Schduler{
		taskChan:make(chan Task,10),
		wg: &sync.WaitGroup{},
	}
	for i:=0; i < 10; i++ {
		go s.Worker()
	}
	return s
}

func (s *Schduler)AddTask(task Task)  {
	s.wg.Add(1)
	s.taskList = append(s.taskList, task)
}

func (s *Schduler)Start()  {
	go func() {
		for _,task := range s.taskList  {
			s.taskChan <- task
		}
	}()
	s.wg.Wait()
}

func (s *Schduler)Worker()  {
	for {
		select {
		case task, ok := <-s.taskChan:
			if ok {
				task()
				s.wg.Done()
			}
		}
	}
}