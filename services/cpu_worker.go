package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUWorker struct {
	threshold float64
	interval  time.Duration
	cooldown  time.Duration

	sendAlert func(string)

	stopChan  chan struct{}
	wg        sync.WaitGroup
	lastAlert time.Time
}

func NewCPUWorker(threshold float64, interval, cooldown time.Duration, send func(string)) *CPUWorker {
	return &CPUWorker{
		threshold: threshold,
		interval:  interval,
		cooldown:  cooldown,
		sendAlert: send,
		stopChan:  make(chan struct{}),
	}
}

func (w *CPUWorker) Name() string {
	return "CPUWorker"
}

func (w *CPUWorker) Start() {
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.check()
			case <-w.stopChan:
				return
			}
		}
	}()
}

func (w *CPUWorker) Stop() {
	close(w.stopChan)
	w.wg.Wait()
}

func (w *CPUWorker) check() {
	percent, _ := cpu.Percent(0, false)

	if percent[0] > w.threshold {
		// cooldown check
		if time.Since(w.lastAlert) < w.cooldown {
			return
		}

		msg := fmt.Sprintf("⚠️ CPU High: %.2f%% (threshold %.0f%%)", percent[0], w.threshold)
		w.sendAlert(msg)

		w.lastAlert = time.Now()
	}
}
