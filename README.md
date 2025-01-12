# dqlite-vip

**Dqlite-vip** is a lightweight and standalone program for creating a highly
available cluster of Linux machines on a local network. Its primary function is
to assign a configurable virtual IP address to one of the machines in the
cluster (i.e., the cluster leader) and to automatically transfer the virtual IP
to another machine in case the current leader fails.

As the name suggests, dqlite-vip is built on [Dqlite](https://dqlite.io/), an
embeddable, and highly available data store powered by the Raft consensus
algorithm.

## Architecture

A typical dqlite-vip cluster consists of 3, 5, or 7 nodes. Since it is built on
Dqlite, which uses the Raft consensus algorithm, all Raft best practices and
considerations apply to dqlite-vip too. Key considerations include:

1. For the cluster to remain operational, the majority of its nodes must be
   available. The majority is defined as `(n/2) + 1`, where `n` is the number of
   nodes in the cluster.
2. Run an odd number of nodes in the cluster. The reason is that a cluster with
   `2n` nodes can tolerate the exact same number of failures as a cluster with
   `2n-1` nodes.
3. Clusters with more than 7 nodes are not recommended due to the additional
   overhead of determining cluster membership and quorum

Once the cluster is up and running, the leader node holds the virtual IP and
broadcasts [Gratuitous ARP](https://wiki.wireshark.org/Gratuitous_ARP) packets
to update the ARP tables of other machines in the network:

![dqlite-vip-cluster](./media/dqlite-vip-cluster.excalidraw.png)

**dqlite-vip** has 4 main components:

- **CLI**: A command-line interface (CLI) to start the program and configure the
  static parameters (e.g., data directory, network interface, etc.).
- **Cluster**: dqlite-vip manages a cluster of Dqlite nodes. The node that is
  elected as the leader will hold the virtual IP address.
- **API**: A REST API that can be used to monitor the status of the cluster and
  configure the virtual IP.
- **VIP**: The VIP component assigns the configured virtual IP address to the
  leader node's network interface and broadcasts Gratuitous ARP packets to
  update the ARP tables of other machines in the network.

## Requirements

1. An amd64 or arm64 Linux machine running a kernel with support for
   [native async I/O](https://man7.org/linux/man-pages/man2/io_setup.2.html).
2. A stable and low-latency network connection between the nodes in the cluster.
3. `CAP_NET_ADMIN` and `CAP_NET_RAW` capabilities are required to run the
   program.

## Additional Considerations

1. The network interface configured with the `cluster-bind` address must retain
   its IP address throughout the program's runtime, even during temporary
   interface downtime.

   > [!NOTE]
   > It's highly recommended to use a static IP address for the network
   > interface.

   > [!TIP]
   > For [netplan](https://netplan.readthedocs.io/en/stable/netplan-yaml), the
   > `ignore-carrier` option can be used to prevent the interface from losing
   > its IP address when the carrier is lost.

   If the interface loses its IP address temporarily, the program won't be able
   to recover until it's restarted.

2. To change the IP addresses of nodes in an existing cluster or to remove
   nodes, you must recreate the cluster and use a clean data directory.

## Quick Start

To start a dqlite-vip cluster, you need to run the program on each node in the
cluster.

The quickest way to get **dqlite-vip** running is to download a pre-built
release binary from the
[releases page](https://github.com/fardjad/dqlite-vip/releases) on GitHub. The
releases are available for _amd64_ and _arm64_ architectures.

```bash
DQLITE_VIP_BINARY=/path/to/dqlite-vip

sudo setcap cap_net_admin,cap_net_raw+ep $DQLITE_VIP_BINARY
sudo install -Dm755 $DQLITE_VIP_BINARY /usr/local/bin/dqlite-vip
```

You can then start the program on the first machine with the following command:

```bash
DATA_DIR=/path/to/data-dir # replace this with the path to the data directory (e.g., /opt/dqlite-vip/data)
NODE_IP="1.2.3.4" # replace this with the node's IP address
IFACE="eth0" # replace this with the network interface to use for the VIP

dqlite-vip start \
  --data-dir "${DATA_DIR}" \
  --bind-cluster "${NODE_IP}:8800" \
  --bind-http "127.0.0.1:9900" \
  --iface "${IFACE}"
```

For additional machines, use the same command but include the `--join` flag to
specify the address of an existing cluster node:

```bash
JOIN_IP="1.2.3.4" # replace this with the IP address of another node in the cluster

dqlite-vip start \
  --data-dir "${DATA_DIR}" \
  --bind-cluster "${NODE_IP}:8800" \
  --bind-http "127.0.0.1:9900" \
  --iface "${IFACE}" \
  --join "${JOIN_IP}:8800"
```

> [!NOTE]
> After the cluster gets formed, the join option will be ignored on subsequent
> runs.

Once all nodes are up and running, you can use the REST API to monitor the state
of the cluster and configure the virtual IP.

Retrieve the cluster status (run this on any node):

```bash
# assuming you have curl and jq installed
curl http://localhost:9900/status | jq
```

Configure the virtual IP (run this on any node):

```bash
VIP="1.2.3.5" # replace with the virtual IP address
SUBNET_PREFIX_LENGTH="24" # replace with the subnet prefix length

curl -XPUT -d '{"vip":"'${VIP}/${SUBNET_PREFIX_LENGTH}'"}' http://localhost:9900/vip
```

After setting the virtual IP, the leader node will set the IP address on the
configured network interface and will start broadcasting Gratuitous ARP packets.

You can then ping the virtual IP from any machine in the network:

```bash
ping -c 1 ${VIP}
```
