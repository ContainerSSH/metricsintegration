package metricsintegration_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/containerssh/geoip"
	"github.com/containerssh/metrics"
	"github.com/containerssh/sshserver"
	"github.com/stretchr/testify/assert"

	"github.com/containerssh/metricsintegration"
)

func TestMetricsReporting(t *testing.T) {
	geoIPProvider, err := geoip.New(geoip.Config{Provider: geoip.DummyProvider})
	if !assert.NoError(t, err) {
		return
	}
	metricsCollector := metrics.New(geoIPProvider)
	backend := &dummyBackendHandler{
		authResponse: sshserver.AuthResponseSuccess,
	}
	handler, err := metricsintegration.NewHandler(
		metrics.Config{
			Enable: true,
		},
		metricsCollector,
		backend,
	)
	if !assert.NoError(t, err) {
		return
	}
	t.Run("auth=successful", func(t *testing.T) {
		testAuthSuccessful(t, handler, metricsCollector)
	})

	t.Run("auth=failed", func(t *testing.T) {
		testAuthFailed(t, backend, handler, metricsCollector)
	})

	t.Run("auth=unavailable", func(t *testing.T) {
		testAuthUnavailable(t, backend, handler, metricsCollector)
	})
}

func testAuthSuccessful(
	t *testing.T,
	handler sshserver.Handler,
	metricsCollector metrics.Collector,
) {
	networkHandler, err := handler.OnNetworkConnection(
		net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 2222,
		},
		sshserver.GenerateConnectionID(),
	)
	if !assert.NoError(t, err) {
		return
	}
	defer networkHandler.OnDisconnect()

	_, err = networkHandler.OnAuthPassword("foo", []byte("bar"))
	if !assert.NoError(t, err) {
		return
	}
	_, err = networkHandler.OnHandshakeSuccess("foo")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameSuccessfulHandshake)[0].Value)
	assert.Equal(t, 0, len(metricsCollector.GetMetric(metricsintegration.MetricNameFailedHandshake)))
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)
	assert.Equal(t, 0, len(metricsCollector.GetMetric(metricsintegration.MetricNameAuthBackendFailure)))
	assert.Equal(t, 0, len(metricsCollector.GetMetric(metricsintegration.MetricNameAuthFailure)))
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameAuthSuccess)[0].Value)

	networkHandler.OnDisconnect()
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(0), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)
}

func testAuthFailed(
	t *testing.T,
	backend *dummyBackendHandler,
	handler sshserver.Handler,
	metricsCollector metrics.Collector,
) {
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(0), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)

	backend.authResponse = sshserver.AuthResponseFailure
	networkHandler, err := handler.OnNetworkConnection(
		net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 2222,
		},
		sshserver.GenerateConnectionID(),
	)
	assert.NoError(t, err)
	response, err := networkHandler.OnAuthPassword("foo", []byte("bar"))
	assert.NoError(t, err)
	assert.Equal(t, sshserver.AuthResponseFailure, response)
	networkHandler.OnHandshakeFailed(fmt.Errorf("failed authentication"))
	assert.Equal(t, float64(2), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameSuccessfulHandshake)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameFailedHandshake)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)
	assert.Equal(t, 0, len(metricsCollector.GetMetric(metricsintegration.MetricNameAuthBackendFailure)))
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameAuthFailure)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameAuthSuccess)[0].Value)

	networkHandler.OnDisconnect()
	assert.Equal(t, float64(2), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(0), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)
}

func testAuthUnavailable(
	t *testing.T,
	backend *dummyBackendHandler,
	handler sshserver.Handler,
	metricsCollector metrics.Collector,
) {
	assert.Equal(t, float64(2), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(0), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)

	backend.authResponse = sshserver.AuthResponseUnavailable
	networkHandler, err := handler.OnNetworkConnection(
		net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 2222,
		},
		sshserver.GenerateConnectionID(),
	)
	assert.NoError(t, err)
	response, err := networkHandler.OnAuthPassword("foo", []byte("bar"))
	assert.NoError(t, err)
	assert.Equal(t, sshserver.AuthResponseUnavailable, response)

	networkHandler.OnHandshakeFailed(fmt.Errorf("auth unavailable"))

	assert.Equal(t, float64(3), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameSuccessfulHandshake)[0].Value)
	assert.Equal(t, float64(2), metricsCollector.GetMetric(metricsintegration.MetricNameFailedHandshake)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameAuthBackendFailure)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameAuthFailure)[0].Value)
	assert.Equal(t, float64(1), metricsCollector.GetMetric(metricsintegration.MetricNameAuthSuccess)[0].Value)

	networkHandler.OnDisconnect()
	assert.Equal(t, float64(3), metricsCollector.GetMetric(metricsintegration.MetricNameConnections)[0].Value)
	assert.Equal(t, float64(0), metricsCollector.GetMetric(metricsintegration.MetricNameCurrentConnections)[0].Value)
}

type dummyBackendHandler struct {
	authResponse sshserver.AuthResponse
}

func (d *dummyBackendHandler) OnReady() error {
	return nil
}

func (d *dummyBackendHandler) OnShutdown(_ context.Context) {

}

func (d *dummyBackendHandler) OnNetworkConnection(
	_ net.TCPAddr,
	_ string,
) (sshserver.NetworkConnectionHandler, error) {
	return d, nil
}

func (d *dummyBackendHandler) OnDisconnect() {
}

func (d *dummyBackendHandler) OnAuthPassword(_ string, _ []byte) (
	response sshserver.AuthResponse,
	reason error,
) {
	return d.authResponse, nil
}

func (d *dummyBackendHandler) OnAuthPubKey(_ string, _ string) (
	response sshserver.AuthResponse,
	reason error,
) {
	return d.authResponse, nil
}

func (d *dummyBackendHandler) OnHandshakeFailed(_ error) {

}

func (d *dummyBackendHandler) OnHandshakeSuccess(_ string) (
	connection sshserver.SSHConnectionHandler,
	failureReason error,
) {
	return d, nil
}

func (d *dummyBackendHandler) OnUnsupportedGlobalRequest(_ uint64, _ string, _ []byte) {

}

func (d *dummyBackendHandler) OnUnsupportedChannel(_ uint64, _ string, _ []byte) {

}

func (d *dummyBackendHandler) OnSessionChannel(
	_ uint64,
	_ []byte,
) (channel sshserver.SessionChannelHandler, failureReason sshserver.ChannelRejection) {
	return d, nil
}

func (d *dummyBackendHandler) OnUnsupportedChannelRequest(_ uint64, _ string, _ []byte) {

}

func (d *dummyBackendHandler) OnFailedDecodeChannelRequest(
	_ uint64,
	_ string,
	_ []byte,
	_ error,
) {

}

func (d *dummyBackendHandler) OnEnvRequest(_ uint64, _ string, _ string) error {
	return fmt.Errorf("env not supported")
}

func (d *dummyBackendHandler) OnPtyRequest(
	_ uint64,
	_ string,
	_ uint32,
	_ uint32,
	_ uint32,
	_ uint32,
	_ []byte,
) error {
	return fmt.Errorf("PTY not supported")
}

func (d *dummyBackendHandler) OnExecRequest(
	_ uint64,
	exec string,
	_ io.Reader,
	stdout io.Writer,
	_ io.Writer,
	exit func(exitStatus sshserver.ExitStatus),
) error {
	go func() {
		_, err := stdout.Write([]byte(fmt.Sprintf("Exec request received: %s", exec)))
		if err != nil {
			exit(2)
		} else {
			exit(0)
		}
	}()
	return nil
}

func (d *dummyBackendHandler) OnShell(
	_ uint64,
	_ io.Reader,
	_ io.Writer,
	_ io.Writer,
	_ func(exitStatus sshserver.ExitStatus),
) error {
	return fmt.Errorf("shell not supported")
}

func (d *dummyBackendHandler) OnSubsystem(
	_ uint64,
	subsystem string,
	_ io.Reader,
	stdout io.Writer,
	_ io.Writer,
	exit func(exitStatus sshserver.ExitStatus),
) error {
	if subsystem != "sftp" {
		return fmt.Errorf("subsystem not supported")
	}
	go func() {
		_, err := stdout.Write([]byte(fmt.Sprintf("Subsystem request received: %s", subsystem)))
		if err != nil {
			exit(2)
		} else {
			exit(0)
		}
	}()
	return nil
}

func (d *dummyBackendHandler) OnSignal(_ uint64, _ string) error {
	return fmt.Errorf("signal not supported")
}

func (d *dummyBackendHandler) OnWindow(
	_ uint64,
	_ uint32,
	_ uint32,
	_ uint32,
	_ uint32,
) error {
	return fmt.Errorf("window changes are not supported")
}
