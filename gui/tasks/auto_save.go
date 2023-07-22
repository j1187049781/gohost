package tasks

import "time"

type AutoSaveTask struct{
	NeedSave chan struct{}
	TaskFunc func()
}

func (t *AutoSaveTask) RunBgTask(){
	go func ()  {
		for {
			<- t.NeedSave
			if t.TaskFunc != nil{
				t.TaskFunc()
			}
			time.Sleep(time.Second * 3)
		}
	}()
}