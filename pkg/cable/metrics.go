package cable

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	submv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
)

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

	connectionEstablishTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "connection_established_timestamp",
			Help: "the Unix timestamp at which the connection established",
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
	prometheus.MustRegister(rxGauge, txGauge, connectionActivationStatus, connectionEstablishTimestamp)
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
	labels := prometheus.Labels{
		cableDriverLabel:    cableDriverName,
		localClusterLabel:   localEndpoint.ClusterID,
		localHostnameLabel:  localEndpoint.Hostname,
		remoteClusterLabel:  remoteEndpoint.ClusterID,
		remoteHostnameLabel: remoteEndpoint.Hostname,
	}
	connectionActivationStatus.With(labels).Set(float64(1))
	connectionEstablishTimestamp.With(labels).Set(float64(time.Now().Unix()))
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
	connectionEstablishTimestamp.Delete(labels)
}
