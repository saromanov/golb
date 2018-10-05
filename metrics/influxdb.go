package metrics

import (
	"github.com/go-kit/kit/metrics/influx"
	influxdb "github.com/influxdata/influxdb/client/v2"
)

var influxDBClient *influx.Influx

// RegisterInfluxDB registers the metrics pusher if this didn't happen yet and creates a InfluxDB Registry instance.
func RegisterInfluxDB() Metrics {
	if influxDBClient == nil {
		influxDBClient = initInfluxDBClient()
	}

	return &simpleMetrics{
		backendReqsCounter:             influxDBClient.NewCounter(influxDBMetricsBackendReqsName),
		backendReqDurationHistogram:    influxDBClient.NewHistogram(influxDBMetricsBackendLatencyName),
		backendRetriesCounter:          influxDBClient.NewCounter(influxDBRetriesTotalName),
	}
}

// initInflux provides initialization of influx db
func initInfluxDB() *influx.Influx{
	return influx.New(
		map[string]string{},
		influxdb.BatchPointsConfig{
			Database:        config.Database,
			RetentionPolicy: config.RetentionPolicy,
		},
		kitlog.LoggerFunc(func(keyvals ...interface{}) error {
			log.Info(keyvals)
			return nil
		}))
}