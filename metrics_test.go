package metrics

import (
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
)

func TestMetricsStoreGraphite(t *testing.T) {
	err := Setup(
		"graphite",
		"localhost",
		2003,
		"",
		time.Hour, // don't flush metrics until Close call
	)
	assert.Nil(t, err, "despite absence of graphite server, metrics setup should end successfully")

	assert.Equal(t, map[string]map[string]interface{}{}, emm.registry.GetAll(), "registry is empty before we start")

	IncrCounter([]string{"happy_routine", "happy_runs_counter"}, 1)
	IncrCounter([]string{"happy_routine", "happy_runs_counter"}, 5)

	UpdateGauge([]string{"happy_routine", "happiness_level"}, 9000)

	UpdateHistogram([]string{"happy_routine", "happiness_hit"}, 35)
	UpdateHistogram([]string{"happy_routine", "happiness_hit"}, 7)

	UpdateTimer([]string{"happy_time"}, time.Minute)
	UpdateTimer([]string{"happy_time"}, time.Second)

	// wait for stats to be dumped
	Close(true)

	assert.Equal(t, int64(6), emm.registry.Get("happy_routine.happy_runs_counter").(metrics.Counter).Count())
	assert.Equal(t, int64(9000), emm.registry.Get("happy_routine.happiness_level").(metrics.Gauge).Value())
	assert.Equal(t, int64(2), emm.registry.Get("happy_routine.happiness_hit").(metrics.Histogram).Count())
	assert.Equal(t, int64(42), emm.registry.Get("happy_routine.happiness_hit").(metrics.Histogram).Sum())
	assert.Equal(t, int64(2), emm.registry.Get("happy_time").(metrics.Timer).Count())
	assert.Equal(t, (time.Minute + time.Second).Nanoseconds(), emm.registry.Get("happy_time").(metrics.Timer).Sum())

	err = Setup(
		"graphite",
		"localhost",
		2003,
		"",
		time.Hour,
	)
	assert.Nil(t, err, "second call should do nothing and shouldn't return error")
	// cleanup
	emm.registry.UnregisterAll()
}
