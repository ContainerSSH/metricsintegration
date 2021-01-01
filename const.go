package metricsintegration

// MetricNameConnections is the number of connections since start.
const MetricNameConnections = "containerssh_ssh_connections"

// MetricHelpConnections is the help text for the number of connections since start.
const MetricHelpConnections = "Number of connections since start"

// MetricNameCurrentConnections is the number of currently open SSH connections.
const MetricNameCurrentConnections = "containerssh_ssh_current_connections"

// MetricHelpCurrentConnections is th ehelp text for the number of currently open SSH connections.
const MetricHelpCurrentConnections = "Current open SSH connections"

// MetricNameSuccessfulHandshake is the number of successful SSH handshakes since start.
const MetricNameSuccessfulHandshake = "containerssh_ssh_handshake_successful"

// MetricHelpSuccessfulHandshake is the help text for the number of successful SSH handshakes since start.
const MetricHelpSuccessfulHandshake = "Successful SSH handshakes since start"

// MetricNameFailedHandshake is the number of failed SSH handshakes since start.
const MetricNameFailedHandshake = "containerssh_ssh_handshake_failed"

// MetricHelpFailedHandshake is the help text for the number of failed SSH handshakes since start.
const MetricHelpFailedHandshake = "Failed SSH handshakes since start"

// MetricNameAuthBackendFailure is the number of request failures to the authentication backend.
const MetricNameAuthBackendFailure = "containerssh_auth_server_failures"

// MetricHelpAuthBackendFailure is the help text for the number of request failures to the authentication backend.
const MetricHelpAuthBackendFailure = "Number of request failures to the authentication backend"

// MetricNameAuthFailure is the number of failed authentications.
const MetricNameAuthFailure = "containerssh_auth_failures"

// MetricHelpAuthFailure is the help text for the number of failed authentications.
const MetricHelpAuthFailure = "Number of failed authentications"

// MetricNameAuthSuccess is the number of successful authentications.
const MetricNameAuthSuccess = "containerssh_auth_success"

// MetricHelpAuthSuccess is the help text for the number of successful authentications.
const MetricHelpAuthSuccess = "Number of successful authentications"
