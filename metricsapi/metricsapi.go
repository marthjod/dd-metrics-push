package metricsapi

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

// MetricsAPI wraps a https://pkg.go.dev/github.com/DataDog/datadog-api-client-go/v2/api/datadogV2#MetricsApi.
type MetricsAPI struct {
	metricsAPIClient *datadogV2.MetricsApi
}

// New returns a ready-to-use MetricsAPI.
func New(ddMetricsAPI *datadogV2.MetricsApi) *MetricsAPI {
	return &MetricsAPI{
		metricsAPIClient: ddMetricsAPI,
	}
}

// Submit wraps https://pkg.go.dev/github.com/DataDog/datadog-api-client-go/v2/api/datadogV2#MetricsApi.SubmitMetrics.
func (m *MetricsAPI) Submit(ctx context.Context, series []datadogV2.MetricSeries) (*http.Response, error) {
	metricPayload := datadogV2.NewMetricPayload(series)
	intakePayloadAccepted, resp, err := m.metricsAPIClient.SubmitMetrics(ctx, *metricPayload)
	if err != nil {
		return nil, err
	}

	// HasErrors() didn't work as expected
	if len(intakePayloadAccepted.GetErrors()) != 0 {
		// TODO: compile errors better?
		return nil, errors.New(strings.Join(intakePayloadAccepted.GetErrors(), ","))
	}

	return resp, nil
}
