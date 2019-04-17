package metrics

import (
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
)

func TestHelpers(t *testing.T) {
	_ = Setup(
		"graphite",
		"localhost",
		2003,
		"",
		time.Hour, // don't flush metrics until Close call
	)

	go Scheduler(time.Hour, MeasureRuntime, []string{"uptime"}, time.Time{})
	GoRuntimeStats([]string{"r"})

	// wait for stats to be dumped
	Close(true)

	// check GoRuntimeStats
	assert.Nil(t, emm.registry.Get("runtime.random_metric"), "random metric is not present")
	assert.NotNil(t, emm.registry.Get("r.runtime.num_goroutines"))
	assert.NotNil(t, emm.registry.Get("r.runtime.heap_alloc"))
	assert.NotNil(t, emm.registry.Get("r.runtime.sys"))
	assert.NotNil(t, emm.registry.Get("r.runtime.pause_total_ns"))
	assert.NotNil(t, emm.registry.Get("r.runtime.num_gc"))
	assert.NotNil(t, emm.registry.Get("r.runtime.heap_released"))
	assert.NotNil(t, emm.registry.Get("r.runtime.heap_objects"))

	// check Scheduler and MeasureRuntime
	// Scheduler problem cases are not tested
	assert.Equal(t, (time.Now().Sub(time.Time{}) / time.Millisecond).Nanoseconds(), emm.registry.Get("uptime").(metrics.Gauge).Value())

	// cleanup
	emm.registry.UnregisterAll()
}
