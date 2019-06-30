package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	//maxTargetsPerFping    = 100
	defaultMinWait         = 10
	startDelayMilliseconds = 1000
)

type WorkerSpec struct {
	period  time.Duration
	count   uint
	minWait uint
}

type Worker struct {
	sync.Mutex
	spec    WorkerSpec
	targets map[TargetSpec]*Target
}

func NewWorker(spec WorkerSpec) *Worker {
	log.Println("New worker (period:", spec.period, ")")
	// initialize defaults
	spec.count = opts.Count
	spec.minWait = defaultMinWait

	// if the period is very short, shorten as well minWait
	if spec.period < 15*time.Second {
		spec.minWait = 1
	}

	// create Worker
	w := Worker{
		spec:    spec,
		targets: make(map[TargetSpec]*Target),
	}

	// TODO: reject if period too small (i.e. < 2 seconds)

	// start main loop
	go w.cycleRun(startDelayMilliseconds * time.Millisecond)

	return &w
}

func (w *Worker) GetWorkerTarget(ts TargetSpec) *Target {
	// TODO: delete unused targets (i.e. if not called for more than 2 periods or so)

	w.Lock()
	defer w.Unlock()
	t, ok := w.targets[ts]
	if !ok {
		t = NewTarget(ts)
		w.targets[ts] = t
	}
	return t
}

func (w *Worker) cycleRun(sleepTime time.Duration) {
	time.Sleep(sleepTime)

	// TODO: only run fping with at most maxTargetsPerFping
	// -> launch multiple go routines

	// schedule the next run
	go w.cycleRun(w.spec.period)

	// prepare fping arguments
	fpingArgs := []string{
		"-q", // quiet
		"-p", // period
		fmt.Sprintf("%.0f", w.spec.period.Seconds()*500/float64(w.spec.count)),
		"-C", // count
		strconv.FormatUint(uint64(w.spec.count), 10),
		"-i", // min-wait
		strconv.FormatUint(uint64(w.spec.minWait), 10),
	}
	for _, t := range w.targets {
		fpingArgs = append(fpingArgs, t.spec.host)
	}

	// start fping
	ctx, cancel := context.WithTimeout(context.Background(), w.spec.period)
	defer cancel()
	fmt.Println("start fping: ", fpingArgs)
	cmd := exec.CommandContext(ctx, opts.Fping, fpingArgs...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		exitErr := err.(*exec.ExitError)
		ws := exitErr.Sys().(syscall.WaitStatus)
		// Note: exit 1 is also returned if any host is unreachable
		if ws.ExitStatus() != 1 {
			fmt.Printf("fping error (exit: %d)", ws.ExitStatus())
			return
		}
	}

	w.addResults(errbuf.String())
}

func (w *Worker) addResults(fpingOutput string) {
	scanner := bufio.NewScanner(strings.NewReader(fpingOutput))
	for scanner.Scan() {
		// Split host and results
		text := strings.SplitN(scanner.Text(), " : ", 2)
		if len(text) != 2 {
			log.Println("Error parsing fping output: ", scanner.Text())
			continue
		}

		// Find target
		host := TargetSpec{host: strings.TrimSpace(text[0])}
		t, ok := w.targets[host]
		if !ok {
			log.Println("Error: fping result for unknown target: ", text[0])
			continue
		}

		// Parse results
		measurements, err := ParseMeasurements(text[1])
		if err != nil {
			log.Println("Error parsing fping output: ", text[1])
			continue
		}

		t.AddMeasurements(measurements)
	}
}
