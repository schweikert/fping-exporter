package main

import (
	"sync"
)

// WorkerSpec and TargetSpec are given by the client and represent respectively
// fping parameters that are common across a group of targets, and the targets
// (hosts to ping). (for WorkerSpec -> see worker.go)

type TargetSpec struct {
	host string
}

// TargetMap maps a WorkerSpec+TargetSpec to a Worker and Target
type TargetMap struct {
	sync.Mutex
	workers map[WorkerSpec]*Worker
}

var tm TargetMap = TargetMap{
	workers: make(map[WorkerSpec]*Worker),
}

func GetTarget(ws WorkerSpec, ts TargetSpec) *Target {
	// retrieve Worker
	tm.Lock()
	w, ok := tm.workers[ws]
	if !ok {
		w = NewWorker(ws)
		tm.workers[ws] = w
	}
	tm.Unlock()

	// retrieve Target
	return w.GetWorkerTarget(ts)
}
