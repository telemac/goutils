package throughput

import (
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	assert := assert.New(t)

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	Formatter := new(log.JSONFormatter)
	log.SetFormatter(Formatter)
	//log.WithField("field", 123456).Info("Info log")

	m := NewMeasurement()
	for i := 0; i < 10; i++ {
		m.AddIterations(1)
		time.Sleep(time.Millisecond * 100)
	}
	r := m.GetThroughput()
	r.Log("10 iterations per secont")
	log.Printf("throughput = %f, iter = %d, duration = %s", r.IterationsPerSecond, r.Iterations, r.Duration)
	assert.InDelta(10, r.IterationsPerSecond, 0.5, "must be about 10 iterations per second")

	// Reset throughput measurment
	m.StartMeasurement()
	time.Sleep(100 * time.Millisecond)
	m.AddIterations(1)
	r = m.GetThroughput()
	assert.InDelta(10, r.IterationsPerSecond, 0.5, "must be about 10 iterations per second")
}
