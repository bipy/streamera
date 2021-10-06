package server

import (
    "fmt"
    "streamera/common"
    "sync"
    "time"
)

type SpeedCounter struct {
    PktCount       int64
    LastTime       int64
    SpeedPerSecond int64
    PingRealTime   int64
    PingPerSecond  int64
    Mutex          sync.RWMutex
}

func NewSpeedCounter() *SpeedCounter {
    counter := &SpeedCounter{
        PktCount:       0,
        LastTime:       time.Now().UnixMicro(),
        SpeedPerSecond: 0,
        PingRealTime:   1000,
        PingPerSecond:  1000,
        Mutex:          sync.RWMutex{},
    }
    return counter
}

func getHumanReadableSpeed(bps int64) string {
    if bps < 1.5e6 {
        return fmt.Sprintf("%d Kbps", bps/1e3)
    }
    return fmt.Sprintf("%.3f Mbps", float64(bps)/1e6)
}

func getHumanReadablePing(ping int64) string {
    if ping >= 1e9 {
        return fmt.Sprintf("999+ ms")
    }
    return fmt.Sprintf("%.3f ms", float64(ping)/float64(time.Millisecond.Microseconds()))
}

func calcSpeed(counter *SpeedCounter) {
    counter.Mutex.RLock()
    lastCount := counter.PktCount
    lastTime := time.Now().UnixMicro()
    counter.Mutex.RUnlock()
    for {
        time.Sleep(time.Second)

        counter.Mutex.Lock()
        curTime := time.Now().UnixMicro()
        curCount := counter.PktCount
        counter.SpeedPerSecond = (curCount - lastCount) * common.ChunkSize * 8 * time.Second.Microseconds() / (curTime - lastTime)
        counter.Mutex.Unlock()
        lastTime, lastCount = curTime, curCount
    }
}

func updatePing(counter *SpeedCounter) {
    for {
        counter.Mutex.Lock()
        counter.PingPerSecond = counter.PingRealTime
        counter.Mutex.Unlock()

        time.Sleep(time.Second)
    }
}
