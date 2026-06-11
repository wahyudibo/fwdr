# fwdr

Have a database, cache, or service that's only reachable inside a private network? Deploy fwdr on a Kubernetes cluster or bastion host with internal access and tunnel to those resources straight from your local machine.

Dead simple TCP forwarder. Listens on a local port and proxies all traffic to a remote host — no config overhead, no dependencies, runs anywhere.

Works with [Kubernetes port-forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) for persistent, self-restarting tunnels — or run it directly on any bastion host.

## Features

- Bidirectional TCP proxying
- Multiple forwarding rules from a single `settings.yaml` — connect to multiple services, each mapped to its own local port
- Graceful shutdown — active connections drain before exit
- Kubernetes-native via Helm chart
- Runs on any VPS or bastion host — drop in a Docker container or a single binary, no dependencies

## Requirements

| Environment | Requirements     |
|-------------|------------------|
| Kubernetes  | Helm 3, kubectl  |
| Docker      | Docker Engine    |
| Binary      | —                |

## Quick Start (Kubernetes)

Let's say you have a Kubernetes cluster and you need to access a database and cache that are only accessible in the internal network. Create a `values.yaml` with your forwarding rules:

```yaml
# my-values.yaml
forwarders:
  - name: my-database
    source: <DATABASE_HOST>:<DATABASE_PORT>
    destinationPort: <DATABASE_DESTINATION_PORT>
  - name: my-cache
    source: <CACHE_HOST>:<CACHE_PORT>
    destinationPort: <CACHE_DESTINATION_PORT>
```

| Field             | Description                                       |
|-------------------|---------------------------------------------------|
| `name`            | Identifier for the connection (logging)           |
| `source`          | Remote `host:port` to forward to                  |
| `destinationPort` | Container port exposed for `kubectl port-forward` |

Deploy using the Helm chart:

```sh
helm install -n <NAMESPACE> fwdr oci://ghcr.io/wahyudibo/fwdr --version <VERSION> -f my-values.yaml
```

Once the pods are running, look up the pod name and start port-forwarding:

```sh
kubectl get pods -n <NAMESPACE>
kubectl -n <NAMESPACE> port-forward pod/<POD_NAME> <DESTINATION_PORT>
```

Your database/cache is now accessible at `localhost:<DESTINATION_PORT>`. To map to a different local port:

```sh
kubectl -n <NAMESPACE> port-forward pod/<POD_NAME> <LOCAL_PORT>:<DESTINATION_PORT>
```

**Upgrade**

```sh
helm upgrade -n <NAMESPACE> fwdr oci://ghcr.io/wahyudibo/fwdr --version <NEW_VERSION> -f my-values.yaml
```

Forwarding rules live in a ConfigMap — the Deployment restarts automatically when the ConfigMap changes.

**Uninstall**

```sh
helm uninstall -n <NAMESPACE> fwdr
```

---

## Docker

Suitable for any VPS or bastion host with Docker installed. Create a `settings.yaml` with your forwarding rules:

```yaml
- name: my-database
  source: <DATABASE_HOST>:<DATABASE_PORT>
  destination_port: <DATABASE_DESTINATION_PORT>

- name: my-cache
  source: <CACHE_HOST>:<CACHE_PORT>
  destination_port: <CACHE_DESTINATION_PORT>
```

Run with ports published to the host (add one `-p` per `destination_port`):

```sh
docker run --rm \
  -v $(pwd)/settings.yaml:/app/settings.yaml \
  -p <LOCAL_PORT>:<DESTINATION_PORT> \
  ghcr.io/wahyudibo/fwdr:latest
```

Also available on Docker Hub:

```sh
docker run --rm \
  -v $(pwd)/settings.yaml:/app/settings.yaml \
  -p <LOCAL_PORT>:<DESTINATION_PORT> \
  wahyudibo/fwdr:latest
```

Run as a persistent background service:

```sh
docker run -d --restart unless-stopped \
  -v $(pwd)/settings.yaml:/app/settings.yaml \
  -p <LOCAL_PORT>:<DESTINATION_PORT> \
  --name fwdr \
  ghcr.io/wahyudibo/fwdr:latest
```

---

## Binary

The simplest option for a VPS or bastion host — single file, no runtime required. Download the latest binary for your platform from the [Releases](https://github.com/wahyudibo/fwdr/releases) page.

```sh
# Linux (amd64)
curl -L https://github.com/wahyudibo/fwdr/releases/latest/download/fwdr-linux-amd64.tar.gz | tar -xz
chmod +x fwdr

# macOS (Apple Silicon)
curl -L https://github.com/wahyudibo/fwdr/releases/latest/download/fwdr-darwin-arm64.tar.gz | tar -xz
chmod +x fwdr
```

Run a forwarding rule:

```sh
./fwdr --name my-database --source <DATABASE_HOST>:<DATABASE_PORT> --destination-port <DATABASE_DESTINATION_PORT>
```

For multiple forwarders, run each as a background process:

```sh
./fwdr --name my-database --source <DATABASE_HOST>:<DATABASE_PORT> --destination-port <DATABASE_DESTINATION_PORT> &
./fwdr --name my-cache --source <CACHE_HOST>:<CACHE_PORT> --destination-port <CACHE_DESTINATION_PORT> &
```

**Running as a systemd service**

To run fwdr at boot, set it up as a systemd service.

1. Copy the binary:

```sh
sudo cp fwdr /usr/local/bin/fwdr
```

2. Copy the service unit and add your forwarding rule to `ExecStart`:

```sh
sudo cp deployment/systemd/fwdr.service /etc/systemd/system/fwdr.service
sudo nano /etc/systemd/system/fwdr.service
```

Edit the `ExecStart` line with your flags:

```ini
ExecStart=/usr/local/bin/fwdr --name my-database --source <DATABASE_HOST>:<DATABASE_PORT> --destination-port <DATABASE_DESTINATION_PORT>
```

For multiple forwarders, create one service file per rule.

3. Enable and start:

```sh
sudo systemctl daemon-reload
sudo systemctl enable fwdr
sudo systemctl start fwdr
```

Check status:

```sh
systemctl status fwdr
```

Follow logs:

```sh
journalctl -u fwdr -f
```

**Available flags**

| Flag                | Default | Description                                  |
|---------------------|---------|----------------------------------------------|
| `--name`            | —       | Connection name (required)                   |
| `--source`          | —       | Remote `host:port` to forward to (required)  |
| `--destination-port`| `8080`  | Local port to listen on                      |
| `--dial-timeout`    | `30s`   | Timeout when connecting to source            |

## License

[MIT](LICENSE)
