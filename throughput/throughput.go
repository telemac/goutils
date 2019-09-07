package throughput

// package throughput is a minimalistic package to measure throughput

import (
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// Measurement measures throughput with StartMeasurement / AddIterations
type Measurement struct {
	Iterations uint64    // Number of iterations done
	StartTime  time.Time // time of measurement startup
}

// Result is the result of a throughput measurement
type Result struct {
	Iterations          uint64        // Number of iterations done
	Duration            time.Duration // measurement duration
	IterationsPerSecond float64
}

// NewMeasurement creates and starts a Throughput measurment
func NewMeasurement() *Measurement {
	var m Measurement
	m.StartMeasurement()
	return &m
}

// StartMeasurement starts the Throughput measurment
func (m *Measurement) StartMeasurement() {
	m.Iterations = 0
	m.StartTime = time.Now()
}

// AddIterations adds iterations to the Throughput measurment
// GetThroughput and AddIterations can be called in different goroutines
func (m *Measurement) AddIterations(nb uint64) {
	atomic.AddUint64(&m.Iterations, nb)
}

// GetThroughput returns the number of iterations per second
// returns Throughput, iterations, duration
// GetThroughput and AddIterations can be called in different goroutines
func (m *Measurement) GetThroughput() Result {
	//duration := time.Now().Sub(t.StartTime)
	var r Result
	r.Duration = time.Since(m.StartTime)
	r.Iterations = atomic.LoadUint64(&m.Iterations)
	r.IterationsPerSecond = float64(r.Iterations) / r.Duration.Seconds()
	return r
}

// Log sneds the output to log, prefixed with msg
func (r Result) Log(msg string) {
	log.Printf("%s : throughput = %f, %d interatios in %s", msg, r.IterationsPerSecond, r.Iterations, r.Duration)
}
