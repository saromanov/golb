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
func RegisterInfluxDB(conf *InfluxConfig) Metrics {
	if influxDBClient == nil {
		influxDBClient = initInfluxDB(conf)
	}

	return &simpleMetrics{
		requestsCounter: &SimpleCounter{
			name: influxDBMetricsRequestsTotal,
		},
		requestsDurationHistogram: &SimpleHistogram{
			name: influxDBMetricsRequestsHistogram,
		},
		requestsGauge: &SimpleGauge{
			name: influxDBMetricsRequestsGauge,
		},
	}
}

// initInflux provides initialization of influx db
func initInfluxDB(conf *InfluxConfig) influxdb.Client {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     conf.GetAddress(),
		Username: "inflixdb",
		Password: "influxdb",
	})
	if err != nil {
		panic(fmt.Sprintf("unable to init influx db: %v", err))
	}

	return c
}

// Write provides writing to influx
func Write(precision string) error {
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  "influxdb",
		Precision: precision,
	})
	if err != nil {
		return err
	}

	return influxDBClient.Write(bp)
}
