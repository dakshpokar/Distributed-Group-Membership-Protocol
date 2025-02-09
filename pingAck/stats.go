package pingAck

import (
	"fmt"
	"time"
)

func NetworkStats() {
	BrLock.RLock()
	if time.Since(StartTime) >= 1*time.Second {
		elapsedTime := time.Since(StartTime).Seconds()
		rxBandwidth := float64(Bytes) / elapsedTime

		fmt.Printf("Bandwidth: %.2f B/s\n", rxBandwidth)
	}
	BrLock.RUnlock()
}
