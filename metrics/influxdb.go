package metrics

import (
	"fmt"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

var influxDBClient influxdb.Client

const (
	influxDBMetricsRequestsTotal     = "golb.requests.total"
	influxDBMetricsRequestsHistogram = "golb.requests.histogram"
	influxDBMetricsRequestsGauge     = "golb.requests.gauge"
)

// RegisterInfluxDB registers the metrics pusher if this didn't happen yet and creates a InfluxDB Registry instance.
func RegisterInfluxDB() Metrics {
	if influxDBClient == nil {
		influxDBClient = initInfluxDB()
	}

	return &simpleMetrics{
		requestsCounter:           &SimpleCounter{},
		requestsDurationHistogram: &SimpleHistogram{},
		requestsGauge:             &SimpleGauge{},
	}
}

// initInflux provides initialization of influx db
func initInfluxDB() influxdb.Client {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "inflixdb",
		Password: "influxdb",
	})
	if err != nil {
		panic(fmt.Sprintf("unable to init influx db: %v", err))
	}

	return c
}
