# WireGuard Cable Driver

[WireGuard](https://www.wireguard.com "WireGuard homepage") is an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography.

Traffic is encrypted and encapsulated in UDP packets.

## Driver design

- WireGuard creates a virtual network device that is accessed via netlink. It appears like any network device and currently has a hardcoded
  name `subwg0`.

- WireGuard identifies peers by their cryptographic public key without the need to exchange shared secrets. The owner of the public key must
  have the corresponding private key to prove identity.

- The driver creates the key pair and adds the public key to the local endpoint so other clusters can connect. Like `ipsec`, the node IP
  address is used as the endpoint udp address of the WireGuard tunnels. A fixed port is used for all endpoints.

- The driver adds routing rules to redirect cross cluster communication through the virtual network device `subwg0`.  (*note: this is
  different from `ipsec`, which intercepts packets at netfilter level.*)

- The driver uses [`wgctrl`](https://github.com/WireGuard/wgctrl-go "WgCtrl github"), a go package that enables control of WireGuard devices
  on multiple platforms. Link creation and removal are done through [`netlink`](https://github.com/vishvananda/netlink "Netlink github").
Currently assuming Linux Kernel WireGuard (`wgtypes.LinuxKernel`).

## Installation

- WireGuard needs to be [installed](https://www.wireguard.com/install "WireGuard installation instructions") on the gateway nodes. For
  example, (Ubuntu < 19.04),

  ```shell
  sudo add-apt-repository ppa:wireguard/wireguard
  sudo apt-get update
  sudo apt-get install linux-headers-`uname -r` -y
  sudo apt-get install wireguard
  ```

- The driver needs to be enabled with

  ```shell
  bin/subctl join --cable-driver wireguard --disable-nat broker-info.subm
  ```

- The default UDP listen port for WireGuard is `5871`. It can be changed by setting the env var `CE_IPSEC_NATTPORT`

## Troubleshooting, limitations

- If you get the following message

  ```text
  Fatal error occurred creating engine: failed to add wireguard device: operation not supported
  ```

  you probably did not install WireGuard on the Gateway node.

- The e2e tests can be run with WireGuard by setting `DEPLOY_ARGS` before calling `make e2e`

  ```shell
  export DEPLOY_ARGS="--deploytool operator --deploytool_submariner_args '--cable-driver=wireguard'"
  ```

- No new `iptables` rules were added, although source NAT needs to be disabled for cross cluster communication. This is similar to disabling
  SNAT when sending cross-cluster traffic between nodes to `submariner-gateway`, so the existing rules should be enough.  **The driver will
fail if the CNI does SNAT before routing to Wireguard** (e.g., failed with Calico, works with Flannel).

## Monitoring

the cabledriver 
The following metrics are exposed currently:
* metrics that exposed per gateway:
    * `wireguard_connected_endpoints` the number of connections.
* metrics that exposed per connection:
    * `wireguard_connection_lifetime` the wireguard connection lifetime in seconds.
    * `wireguard_tx_bytes` Bytes transmitted for the connection.
    * `wireguard_rx_bytes` Bytes received for the connection.
### example
for example we have 2 connected clusters , `cluster1` and `cluster2`.
```
Showing information for cluster "kubernetes-admin@cluster.local":
GATEWAY                         CLUSTER                 REMOTE IP       CABLE DRIVER            SUBNETS                                 STATUS
omer-sub2-vm-worker2            cluster2                10.243.64.7     wireguard               10.234.0.0/18, 10.234.64.0/18           connected
```
when we curl to `cluster1` gateway node with `curl 10.243.64.6:8080/metrics` we can observe the metrics:
```
# HELP wireguard_connected_endpoints wireguard connected endpoints
# TYPE wireguard_connected_endpoints gauge
wireguard_connected_endpoints 1
# HELP wireguard_connection_lifetime wireguard connection lifetime in seconds
# TYPE wireguard_connection_lifetime gauge
wireguard_connection_lifetime{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker2",dst_PrivateIP="10.243.64.7",dst_PublicIP="",dst_clusterID="cluster2"} 2488
# HELP wireguard_rx_bytes Bytes received
# TYPE wireguard_rx_bytes gauge
wireguard_rx_bytes{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker2",dst_PrivateIP="10.243.64.7",dst_PublicIP="",dst_clusterID="cluster2"} 7716
# HELP wireguard_tx_bytes Bytes transmitted
# TYPE wireguard_tx_bytes gauge
wireguard_tx_bytes{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker2",dst_PrivateIP="10.243.64.7",dst_PublicIP="",dst_clusterID="cluster2"} 6048
```
after adding another cluster called `cluster3`:
```
Showing information for cluster "kubernetes-admin@cluster.local":
GATEWAY                         CLUSTER                 REMOTE IP       CABLE DRIVER        SUBNETS                                 STATUS
omer-sub2-vm-worker2            cluster2                10.243.64.7     wireguard           10.234.0.0/18, 10.234.64.0/18           connected
omer-sub2-vm-worker3            cluster3                10.243.64.9     wireguard           11.235.0.0/18, 11.235.64.0/18           connected
```
when we curl again to `cluster1` gateway node with `curl 10.243.64.6:8080/metrics` we can observe the metrics:
```
# HELP wireguard_connected_endpoints wireguard connected endpoints
# TYPE wireguard_connected_endpoints gauge
wireguard_connected_endpoints 2
# HELP wireguard_connection_lifetime wireguard connection lifetime in seconds
# TYPE wireguard_connection_lifetime gauge
wireguard_connection_lifetime{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker2",dst_PrivateIP="10.243.64.7",dst_PublicIP="",dst_clusterID="cluster2"} 2904
wireguard_connection_lifetime{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker3",dst_PrivateIP="10.243.64.9",dst_PublicIP="",dst_clusterID="cluster3"} 83
# HELP wireguard_rx_bytes Bytes received
# TYPE wireguard_rx_bytes gauge
wireguard_rx_bytes{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker2",dst_PrivateIP="10.243.64.7",dst_PublicIP="",dst_clusterID="cluster2"} 8896
wireguard_rx_bytes{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker3",dst_PrivateIP="10.243.64.9",dst_PublicIP="",dst_clusterID="cluster3"} 308
# HELP wireguard_tx_bytes Bytes transmitted
# TYPE wireguard_tx_bytes gauge
wireguard_tx_bytes{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker2",dst_PrivateIP="10.243.64.7",dst_PublicIP="",dst_clusterID="cluster2"} 6996
wireguard_tx_bytes{Backend="wireguard",dst_EndPoint_hostname="omer-sub2-vm-worker3",dst_PrivateIP="10.243.64.9",dst_PublicIP="",dst_clusterID="cluster3"} 400
```
### known issues
- if one removes manually a peer from the `submariner` wireguard interface the metrics that exposed for that connection Won't delete.
