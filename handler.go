package metricsintegration

import (
	"context"
	"net"
	"sync"

	"github.com/containerssh/metrics"
	"github.com/containerssh/sshserver"
)

type metricsHandler struct {
	backend                   sshserver.Handler
	metricsCollector          metrics.Collector
	connectionsMetric         metrics.SimpleGeoCounter
	handshakeSuccessfulMetric metrics.SimpleGeoCounter
	handshakeFailedMetric     metrics.SimpleGeoCounter
	currentConnectionsMetric  metrics.SimpleGeoGauge
	authBackendFailureMetric  metrics.SimpleCounter
	authFailureMetric         metrics.SimpleGeoCounter
	authSuccessMetric         metrics.SimpleGeoCounter
}

func (m *metricsHandler) OnReady() error {
	return m.backend.OnReady()
}

func (m *metricsHandler) OnShutdown(shutdownContext context.Context) {
	m.backend.OnShutdown(shutdownContext)
}

func (m *metricsHandler) OnNetworkConnection(
	client net.TCPAddr,
	connectionID string,
) (sshserver.NetworkConnectionHandler, error) {

	networkBackend, err := m.backend.OnNetworkConnection(client, connectionID)
	if err != nil {
		return networkBackend, err
	}
	m.connectionsMetric.Increment(client.IP)
	m.currentConnectionsMetric.Increment(client.IP)
	return &metricsNetworkHandler{
		client:  client,
		backend: networkBackend,
		handler: m,
		lock:    &sync.Mutex{},
	}, nil
}

type metricsNetworkHandler struct {
	backend      sshserver.NetworkConnectionHandler
	client       net.TCPAddr
	handler      *metricsHandler
	lock         *sync.Mutex
	disconnected bool
}

func (m *metricsNetworkHandler) OnShutdown(shutdownContext context.Context) {
	m.backend.OnShutdown(shutdownContext)
}

func (m *metricsNetworkHandler) OnAuthPassword(username string, password []byte) (
	response sshserver.AuthResponse,
	reason error,
) {
	label := metrics.Label("authtype", "password")
	response, reason = m.backend.OnAuthPassword(username, password)
	switch response {
	case sshserver.AuthResponseSuccess:
		m.handler.authSuccessMetric.Increment(m.client.IP, label)
	case sshserver.AuthResponseFailure:
		m.handler.authFailureMetric.Increment(m.client.IP, label)
	case sshserver.AuthResponseUnavailable:
		m.handler.authBackendFailureMetric.Increment(label)
	}
	return response, reason
}

func (m *metricsNetworkHandler) OnAuthPubKey(username string, pubKey string) (
	response sshserver.AuthResponse,
	reason error,
) {
	label := metrics.Label("authtype", "pubkey")
	response, reason = m.backend.OnAuthPubKey(username, pubKey)
	switch response {
	case sshserver.AuthResponseSuccess:
		m.handler.authSuccessMetric.Increment(m.client.IP, label)
	case sshserver.AuthResponseFailure:
		m.handler.authFailureMetric.Increment(m.client.IP, label)
	case sshserver.AuthResponseUnavailable:
		m.handler.authBackendFailureMetric.Increment(label)
	}
	return response, reason
}

func (m *metricsNetworkHandler) OnAuthKeyboardInteractive(
	user string,
	challenge func(
		instruction string,
		questions sshserver.KeyboardInteractiveQuestions,
	) (answers sshserver.KeyboardInteractiveAnswers, err error),
) (response sshserver.AuthResponse, reason error) {
	label := metrics.Label("authtype", "keyboard-interactive")
	response, reason = m.backend.OnAuthKeyboardInteractive(user, challenge)
	switch response {
	case sshserver.AuthResponseSuccess:
		m.handler.authSuccessMetric.Increment(m.client.IP, label)
	case sshserver.AuthResponseFailure:
		m.handler.authFailureMetric.Increment(m.client.IP, label)
	case sshserver.AuthResponseUnavailable:
		m.handler.authBackendFailureMetric.Increment(label)
	}
	return response, reason
}

func (m *metricsNetworkHandler) OnHandshakeFailed(reason error) {
	m.handler.handshakeFailedMetric.Increment(m.client.IP)
	m.backend.OnHandshakeFailed(reason)
}

func (m *metricsNetworkHandler) OnHandshakeSuccess(username string) (
	connection sshserver.SSHConnectionHandler,
	failureReason error,
) {
	connectionHandler, failureReason := m.backend.OnHandshakeSuccess(username)
	if failureReason != nil {
		m.handler.handshakeFailedMetric.Increment(m.client.IP)
		return connectionHandler, failureReason
	}
	m.handler.handshakeSuccessfulMetric.Increment(m.client.IP)
	return connectionHandler, failureReason
}

func (m *metricsNetworkHandler) OnDisconnect() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.disconnected {
		m.handler.currentConnectionsMetric.Decrement(m.client.IP)
		m.disconnected = true
	}
	m.backend.OnDisconnect()
}
