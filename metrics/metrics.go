package metrics

// Metrics provides main interface for metrics
type Metrics interface {
	RequestsCounter() Counter
	RequestsDurationHistogram() Histogram
	RequestsGauge() Gauge
}

// Counter describes a metric that accumulates values monotonically.
// An example of a counter is the number of received HTTP requests.
type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
}

// Gauge describes a metric that takes specific values over time.
// An example of a gauge is the current depth of a job queue.
type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
}

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets. An example of a histogram is
// HTTP request latencies.
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

// simpleMetrics implements metrics for golb
type simpleMetrics struct {
	requestsCounter           Counter
	requestsDurationHistogram Histogram
	requestsGauge             Gauge
}

func (sm *simpleMetrics) RequestsCounter() Counter {
	return sm.requestsCounter
}

func (sm *simpleMetrics) RequestsDurationHistogram() Histogram {
	return sm.requestsDurationHistogram
}

func (sm *simpleMetrics) RequestsGauge() Gauge {
	return sm.requestsGauge
}

type observeFunc func(name string, lvs LabelValues, value float64)

// SimpleCounter is an Influx counter. Observations are forwarded to an Influx
// object, and aggregated (summed) per timeseries.
type SimpleCounter struct {
	name string
	lvs  LabelValues
	obs  observeFunc
}

// With implements metrics.Counter.
func (c *SimpleCounter) With(labelValues ...string) Counter {
	return &SimpleCounter{
		name: c.name,
		lvs:  c.lvs.With(labelValues...),
		obs:  c.obs,
	}
}

// Add implements metrics.Counter.
func (c *SimpleCounter) Add(delta float64) {
	c.obs(c.name, c.lvs, delta)
}

// SimpleHistogram is an Influx histrogram. Observations are aggregated into a
// generic.Histogram and emitted as per-quantile gauges to the Influx server.
type SimpleHistogram struct {
	name string
	lvs  LabelValues
	obs  observeFunc
}

// With implements metrics.Histogram.
func (h *SimpleHistogram) With(labelValues ...string) Histogram {
	return &SimpleHistogram{
		name: h.name,
		lvs:  h.lvs.With(labelValues...),
		obs:  h.obs,
	}
}

// Observe implements metrics.Histogram.
func (h *SimpleHistogram) Observe(value float64) {
	h.obs(h.name, h.lvs, value)
}

// SimpleGauge is an Influx gauge. Observations are forwarded to a Dogstatsd
// object, and aggregated (the last observation selected) per timeseries.
type SimpleGauge struct {
	name string
	lvs  LabelValues
	obs  observeFunc
	add  observeFunc
}

// With implements metrics.Gauge.
func (g *SimpleGauge) With(labelValues ...string) Gauge {
	return &SimpleGauge{
		name: g.name,
		lvs:  g.lvs.With(labelValues...),
		obs:  g.obs,
		add:  g.add,
	}
}

// Set implements metrics.SimpleGauge.
func (g *SimpleGauge) Set(value float64) {
	g.obs(g.name, g.lvs, value)
}

// Add implements metrics.SimpleGauge.
func (g *SimpleGauge) Add(delta float64) {
	g.add(g.name, g.lvs, delta)
}
