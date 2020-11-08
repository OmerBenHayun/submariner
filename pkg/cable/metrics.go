package cable

import (
	"github.com/prometheus/client_golang/prometheus"

	submv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
)

/*
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

// fixme find a better/descriptive name metric for that.adding
var ConnectionActivationStatus = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "connection_activation_status",
		Help: "connection is connected/disconnected.represented as 1/0 respectively",
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
	prometheus.MustRegister(ConnectionActivationStatus, ConnectionUptimeDurationSeconds,
		ConnectionTxBytes, ConnectionRxBytes)
}
*/

const (
	cableDriverLabel    = "cable_driver"
	localClusterLabel   = "local_cluster"
	localHostnameLabel  = "local_hostname"
	remoteClusterLabel  = "remote_cluster"
	remoteHostnameLabel = "remote_hostname"
)

var (
	// The following metrics are gauges because we want to set the absolute value
	// RX/TX metrics
	rxGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_rx_bytes",
			Help: "Count of bytes received (by cable driver and cable)",
		},
		[]string{
			cableDriverLabel,
			localClusterLabel,
			localHostnameLabel,
			remoteClusterLabel,
			remoteHostnameLabel,
		},
	)
	txGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_tx_bytes",
			Help: "Count of bytes transmitted (by cable driver and cable)",
		},
		[]string{
			cableDriverLabel,
			localClusterLabel,
			localHostnameLabel,
			remoteClusterLabel,
			remoteHostnameLabel,
		},
	)
	connectionActivationStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "connection_activation_status",
			Help: "connection is connected/disconnected.represented as 1/0 respectively (by cable driver and cable)",
		},
		[]string{
			cableDriverLabel,
			localClusterLabel,
			localHostnameLabel,
			remoteClusterLabel,
			remoteHostnameLabel,
		},
	)
)

func init() {
	prometheus.MustRegister(rxGauge, txGauge, connectionActivationStatus)
}

func RecordRxBytes(cableDriverName string, localEndpoint, remoteEndpoint *submv1.EndpointSpec, bytes int) {
	rxGauge.With(prometheus.Labels{
		cableDriverLabel:    cableDriverName,
		localClusterLabel:   localEndpoint.ClusterID,
		localHostnameLabel:  localEndpoint.Hostname,
		remoteClusterLabel:  remoteEndpoint.ClusterID,
		remoteHostnameLabel: remoteEndpoint.Hostname,
	}).Set(float64(bytes))
}

func RecordTxBytes(cableDriverName string, localEndpoint, remoteEndpoint *submv1.EndpointSpec, bytes int) {
	txGauge.With(prometheus.Labels{
		cableDriverLabel:    cableDriverName,
		localClusterLabel:   localEndpoint.ClusterID,
		localHostnameLabel:  localEndpoint.Hostname,
		remoteClusterLabel:  remoteEndpoint.ClusterID,
		remoteHostnameLabel: remoteEndpoint.Hostname,
	}).Set(float64(bytes))
}

func RecordConnectionStatusActive(cableDriverName string, localEndpoint, remoteEndpoint *submv1.EndpointSpec) {
	connectionActivationStatus.With(prometheus.Labels{
		cableDriverLabel:    cableDriverName,
		localClusterLabel:   localEndpoint.ClusterID,
		localHostnameLabel:  localEndpoint.Hostname,
		remoteClusterLabel:  remoteEndpoint.ClusterID,
		remoteHostnameLabel: remoteEndpoint.Hostname,
	}).Set(float64(1))
}

func RecordConnectionStatusInactive(cableDriverName string, localEndpoint, remoteEndpoint *submv1.EndpointSpec) {
	labels := prometheus.Labels{
		cableDriverLabel:    cableDriverName,
		localClusterLabel:   localEndpoint.ClusterID,
		localHostnameLabel:  localEndpoint.Hostname,
		remoteClusterLabel:  remoteEndpoint.ClusterID,
		remoteHostnameLabel: remoteEndpoint.Hostname,
	}
	connectionActivationStatus.With(labels).Set(float64(0))
	txGauge.Delete(labels)
	rxGauge.Delete(labels)
}
