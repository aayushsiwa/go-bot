package services

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func StartCPUMonitor(sendAlert func(string)) {
	go func() {
		for {
			percent, _ := cpu.Percent(0, false)

			if percent[0] > 80 {
				sendAlert(
					"⚠️ High CPU Usage: " +
						fmt.Sprintf("%.2f%%", percent[0]),
				)
			}

			time.Sleep(30 * time.Second)
		}
	}()
}
