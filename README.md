[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.github.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Metrics Integration Library</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/metricsintegration?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/library-template)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/metricsintegration?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/library-template/)

This library integrates the [metrics service](https://github.com/containerssh/metrics) with the [sshserver library](https://github.com/containerssh/sshserver).

<p align="center"><strong>⚠⚠⚠ Warning: This is a developer documentation. ⚠⚠⚠</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.io">containerssh.io</a>.</p>

## Using this library

This library is intended as an overlay/proxy for a handler for the [sshserver library](https://github.com/containerssh/sshserver) "handler". It can be injected transparently to collect the following metrics:

- `containerssh_ssh_connections`
- `containerssh_ssh_handshake_successful`
- `containerssh_ssh_handshake_failed`
- `containerssh_ssh_current_connections`
- `containerssh_auth_server_failures`
- `containerssh_auth_failures`
- `containerssh_auth_success`

The handler can be instantiated with the following method:

```go
handler, err := metricsintegration.New(
    config,
    metricsCollector,
    backend,
)
```

- `config` is a configuration structure from the [metrics library](https://github.com/containerssh/metrics). This is used to bypass the metrics integration backend if metrics are disabled.
- `metricsCollector` is the metrics collector from the [metrics library](https://github.com/containerssh/metrics).
- `backend` is an SSH server backend from the [sserver library](https://github.com/containerssh/sshserver).
