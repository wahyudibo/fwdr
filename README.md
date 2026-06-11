# fwdr

Dead simple TCP forwarder. Listens on a local port and proxies all traffic to a remote host — no config overhead, no dependencies, runs anywhere.

Built for [Kubernetes port-forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) workflows where you need persistent, self-restarting tunnels to remote services.

## Features

- Bidirectional TCP proxying
- Multiple forwarding rules from a single `settings.yaml`
- Graceful shutdown — active connections drain before exit
- Kubernetes-native via Helm chart
- Docker-ready for non-Kubernetes environments

## Requirements

| Environment | Requirements     |
|-------------|------------------|
| Kubernetes  | Helm 3, kubectl  |
| Docker      | Docker Engine    |
| Binary      | —                |

## Configuration

Create a `settings.yaml` with your forwarding rules:

```yaml
- name: my-database
  source: db.internal:5432
  destination_port: 5432

- name: my-cache
  source: cache.internal:6379
  destination_port: 6379
```

| Field              | Description                            |
|--------------------|----------------------------------------|
| `name`             | Identifier for the connection (logging)|
| `source`           | Remote `host:port` to forward to       |
| `destination_port` | Local port to listen on                |

## Installation

### Kubernetes (Helm)

Install directly from the OCI registry — no repo clone needed.

```sh
helm install fwdr oci://ghcr.io/wahyudibo/fwdr --version 1.0.0 \
  --set forwarders[0].name=my-database \
  --set forwarders[0].source=db.internal:5432 \
  --set forwarders[0].destinationPort=5432
```

For multiple rules, use a `values.yaml` override file:

```yaml
# my-values.yaml
forwarders:
  - name: my-database
    source: db.internal:5432
    destinationPort: 5432
  - name: my-cache
    source: cache.internal:6379
    destinationPort: 6379
```

```sh
helm install fwdr oci://ghcr.io/wahyudibo/fwdr --version 1.0.0 -f my-values.yaml
```

Forwarding rules are stored in a ConfigMap and mounted into the pod. The Deployment automatically restarts when the ConfigMap changes.

**Upgrade**

```sh
helm upgrade fwdr oci://ghcr.io/wahyudibo/fwdr --version 1.1.0 -f my-values.yaml
```

**Uninstall**

```sh
helm uninstall fwdr
```

---

### Docker

```sh
docker run --rm \
  -v $(pwd)/settings.yaml:/app/settings.yaml \
  ghcr.io/wahyudibo/fwdr:latest
```

Available on Docker Hub too:

```sh
docker run --rm \
  -v $(pwd)/settings.yaml:/app/settings.yaml \
  wahyudibo/fwdr:latest
```

---

### Binary

Download the latest binary for your platform from the [Releases](https://github.com/wahyudibo/fwdr/releases) page.

```sh
# Linux (amd64)
curl -L https://github.com/wahyudibo/fwdr/releases/latest/download/fwdr-linux-amd64 -o fwdr
chmod +x fwdr

# macOS (Apple Silicon)
curl -L https://github.com/wahyudibo/fwdr/releases/latest/download/fwdr-darwin-arm64 -o fwdr
chmod +x fwdr
```

Run a single forwarding rule:

```sh
./fwdr --name my-database --source db.internal:5432 --destination-port 5432
```

**Available flags**

| Flag                | Default | Description                                  |
|---------------------|---------|----------------------------------------------|
| `--name`            | —       | Connection name (required)                   |
| `--source`          | —       | Remote `host:port` to forward to (required)  |
| `--destination-port`| `8080`  | Local port to listen on (required)           |
| `--dial-timeout`    | `30s`   | Timeout when connecting to source            |

## License

[MIT](LICENSE)
