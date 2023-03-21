package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/marthjod/dd-metrics-push/metricsapi"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

const (
	// EnvDDClientAPIKey is the env var key for the DD API key to use.
	EnvDDClientAPIKey = "DD_CLIENT_API_KEY"
	// EnvDDClientAppKey is the env var key for the DD app key to use.
	EnvDDClientAppKey = "DD_CLIENT_APP_KEY"
	// DatadogSiteHost overrides the default configuration's Host field for API calls.
	DatadogSiteHost = "datadoghq.eu"
)

func main() {
	var (
		metricName = flag.String("metric-name", "testing.ignoreme.localsubmit", "Metric name")
		metricTags = flag.String("metric-tags", "testing,ignoreme", "Metric tags (comma-separated)")
	)
	flag.Parse()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	// set up DD connection
	if os.Getenv(EnvDDClientAPIKey) == "" || os.Getenv(EnvDDClientAppKey) == "" {
		log.Fatal("need env vars set: " + strings.Join([]string{EnvDDClientAPIKey, EnvDDClientAppKey}, ", "))
	}

	// https://github.com/DataDog/datadog-api-client-go
	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: os.Getenv(EnvDDClientAPIKey),
			},
			"appKeyAuth": {
				Key: os.Getenv(EnvDDClientAppKey),
			},
		},
	)

	// set up "submitter"
	conf := datadog.NewConfiguration()
	conf.Host = DatadogSiteHost
	metricsAPI := metricsapi.New(datadogV2.NewMetricsApi(datadog.NewAPIClient(conf)))

	// add series and submit them here, potentially in a loop
	// example gauge timeseries
	tags := strings.Split(*metricTags, ",")
	series := []datadogV2.MetricSeries{
		{
			Metric: *metricName,
			Tags:   tags,
			Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
			Points: []datadogV2.MetricPoint{
				{
					Timestamp: datadog.PtrInt64(time.Now().Unix()),
					Value:     datadog.PtrFloat64(rand.NormFloat64()),
				},
			},
		},
	}
	_, err := metricsAPI.Submit(ctx, series)
	if err != nil {
		log.Fatal(err)
	}
}
