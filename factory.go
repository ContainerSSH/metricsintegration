package metricsintegration

import (
	"github.com/containerssh/metrics"
	"github.com/containerssh/sshserver"
)

func NewHandler(
	config metrics.Config,
	metricsCollector metrics.Collector,
	backend sshserver.Handler,
) (sshserver.Handler, error) {
	if !config.Enable {
		return backend, nil
	}

	connectionsMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameConnections,
		"connections",
		MetricHelpConnections,
	)
	currentConnectionsMetric := metricsCollector.MustCreateGaugeGeo(
		MetricNameCurrentConnections,
		"connections",
		MetricHelpCurrentConnections,
	)

	authSuccessMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameAuthSuccess,
		"attempts",
		MetricHelpAuthSuccess,
	)
	authFailureMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameAuthFailure,
		"attempts",
		MetricHelpAuthFailure,
	)
	authBackendFailureMetric := metricsCollector.MustCreateCounter(
		MetricNameAuthBackendFailure,
		"attempts",
		MetricHelpAuthBackendFailure,
	)

	handshakeSuccessfulMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameSuccessfulHandshake,
		"handshakes",
		MetricHelpSuccessfulHandshake,
	)
	handshakeFailedMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameFailedHandshake,
		"handshakes",
		MetricHelpFailedHandshake,
	)

	return &metricsHandler{
		backend:                   backend,
		metricsCollector:          metricsCollector,
		connectionsMetric:         connectionsMetric,
		handshakeSuccessfulMetric: handshakeSuccessfulMetric,
		handshakeFailedMetric:     handshakeFailedMetric,
		currentConnectionsMetric:  currentConnectionsMetric,
		authBackendFailureMetric:  authBackendFailureMetric,
		authSuccessMetric:         authSuccessMetric,
		authFailureMetric:         authFailureMetric,
	}, nil
}
