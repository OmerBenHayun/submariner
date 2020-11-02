package cable

import (
	"github.com/prometheus/client_golang/prometheus"
)

var connectionLabels = []string{
	// destination clusterID
	"clusterID",
	// destination Endpoint hostname
	"hostname",
	// destination PrivateIP
	"privateIP",
	// destination PublicIP
	"publicIP",
	// cable driver name
	"cable_driver",
}

//todo: add wireguard relevat metrics for connectionActivationStatus
var connectionActivationStatus = prometheus.NewGaugeVec( //fixme find a better/descriptive name metric for that.adding
	prometheus.GaugeOpts{
		Name: "connection_activation_status",
		Help: "connection is connected/not connected.represented as 1/0 respectively",
	},
	connectionLabels,
)

var ConnectionUptimeDurationSeconds = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "connection_uptime_duration_seconds",
		Help: "connection uptime duration in seconds",
	},
	connectionLabels,
)

var ConnectionTxBytes = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "connection_tx_bytes",
		Help: "Bytes transmitted to the connection",
	},
	connectionLabels,
)
var ConnectionRxBytes = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "connection_rx_bytes",
		Help: "Bytes received from the connection",
	},
	connectionLabels,
)

func init() {
	prometheus.MustRegister(connectionActivationStatus, ConnectionUptimeDurationSeconds, ConnectionTxBytes, ConnectionRxBytes)
}
