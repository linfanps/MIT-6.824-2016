package mapreduce

import "fmt"
import "sync"

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	var wg sync.WaitGroup
	for i := 0; i < ntasks; i++ {
		wg.Add(1)
		go func(taskNum int) {
			for {
				readyWorker := <-mr.registerChannel
				taskArgs := &DoTaskArgs{
					JobName:       mr.jobName,
					File:          mr.files[taskNum],
					Phase:         phase,
					TaskNumber:    taskNum,
					NumOtherPhase: nios,
				}
				ok := call(readyWorker, "Worker.DoTask", taskArgs, new(struct{}))
				if ok == false {
					fmt.Printf("DoTask: RPC %s do task error\n", readyWorker)
				} else {
					go func() { mr.registerChannel <- readyWorker }()
					wg.Done()
					break
				}
			}
		}(i)
	}
	wg.Wait()

	fmt.Printf("Schedule: %v phase done\n", phase)
}
