# OpenTelemetry Collector Setup for Metrifuge

This document describes how to set up and use the OpenTelemetry Collector to receive logs and metrics from the Metrifuge application.

## Overview

The OpenTelemetry Collector is configured to:
- Receive OTLP metrics and logs via gRPC (port 4317) and HTTP (port 4318)
- Process data with batching and memory limiting
- Export to multiple backends (logging, file, Prometheus)

## Quick Start

### Docker Compose Deployment

The easiest way to run the OpenTelemetry Collector locally is with Docker Compose:

```bash
# Start the collector
docker-compose -f docker-compose.otel.yaml up -d

# View logs
docker-compose -f docker-compose.otel.yaml logs -f otel-collector

# Stop the collector
docker-compose -f docker-compose.otel.yaml down
```

### Kubernetes Deployment

To deploy the OpenTelemetry Collector in Kubernetes:

```bash
# Apply the collector configuration and deployment
kubectl apply -f k8s/otel-collector-deployment.yaml

# Verify the deployment
kubectl get pods -l app=otel-collector
kubectl get svc otel-collector

# View collector logs
kubectl logs -l app=otel-collector -f
```

## Configuration

### Collector Configuration

The collector is configured via `otel-collector-config.yaml` with the following components:

#### Receivers
- **OTLP gRPC** (port 4317): Receives metrics and logs from Metrifuge
- **OTLP HTTP** (port 4318): Alternative HTTP endpoint

#### Processors
- **Batch Processor**: Groups telemetry data for efficient export
  - Timeout: 10 seconds
  - Batch size: 1024 (max 2048)
- **Memory Limiter**: Prevents OOM issues
  - Limit: 512 MiB
  - Spike limit: 128 MiB
- **Resource Processor**: Adds common attributes to all telemetry

#### Exporters
- **Logging Exporter**: Outputs to console (useful for debugging)
- **File Exporters**: Writes metrics and logs to `/var/log/otel/`
  - Rotation: 100 MB max, 7 days retention, 3 backups
- **Prometheus Exporter** (port 8889): Exposes metrics for Prometheus scraping

### Metrifuge Configuration

To send data to the OpenTelemetry Collector, configure Metrifuge exporters:

#### Metrics Export

The Metrifuge application already includes OTLP metric export support. Ensure exporters are configured in your Kubernetes CRDs:

```yaml
apiVersion: metrifuge.devon-caron.dev/v1
kind: Exporter
metadata:
  name: otel-metrics-exporter
spec:
  destinationType: otel_metric_exporter
  endpoint: otel-collector:4317  # In K8s
  # or localhost:4317 for local development
```

#### Logs Export (New)

Log export support has been added via the new `otel_log_exporter_client` package. Configure log exporters:

```yaml
apiVersion: metrifuge.devon-caron.dev/v1
kind: Exporter
metadata:
  name: otel-logs-exporter
spec:
  destinationType: otel_log_exporter
  endpoint: otel-collector:4317
```

## Accessing Telemetry Data

### Prometheus Metrics

The collector exposes metrics in Prometheus format on port 8889:

```bash
# Local
curl http://localhost:8889/metrics

# Kubernetes
kubectl port-forward svc/otel-collector 8889:8889
curl http://localhost:8889/metrics
```

### File Exports

Logs and metrics are written to JSON files:

**Docker:**
```bash
# Access the volume
docker exec -it metrifuge-otel-collector cat /var/log/otel/metrics.json
docker exec -it metrifuge-otel-collector cat /var/log/otel/logs.json
```

**Kubernetes:**
```bash
# Get pod name
POD=$(kubectl get pods -l app=otel-collector -o jsonpath='{.items[0].metadata.name}')

# View metrics file
kubectl exec $POD -- cat /var/log/otel/metrics.json

# View logs file
kubectl exec $POD -- cat /var/log/otel/logs.json
```

### Console Output

View real-time telemetry data via collector logs:

**Docker:**
```bash
docker-compose -f docker-compose.otel.yaml logs -f otel-collector
```

**Kubernetes:**
```bash
kubectl logs -l app=otel-collector -f
```

## Architecture

```
┌─────────────────┐
│   Metrifuge     │
│   Application   │
└────────┬────────┘
         │ OTLP gRPC/HTTP
         │ (metrics & logs)
         ▼
┌─────────────────────┐
│ OpenTelemetry       │
│ Collector           │
│                     │
│ Receivers → Process │
│ ors → Exporters     │
└───────┬─────────────┘
        │
        ├─────► Console (logging exporter)
        │
        ├─────► Files (/var/log/otel/*.json)
        │
        ├─────► Prometheus endpoint (:8889)
        │
        └─────► External backends (configurable)
```

## Ports

| Port | Protocol | Purpose |
|------|----------|---------|
| 4317 | gRPC | OTLP receiver (metrics & logs) |
| 4318 | HTTP | OTLP receiver (alternative) |
| 8889 | HTTP | Prometheus metrics exporter |
| 8888 | HTTP | Collector internal metrics |
| 13133 | HTTP | Health check endpoint |

## Troubleshooting

### Collector Not Receiving Data

1. Check that Metrifuge can reach the collector:
   ```bash
   # In Kubernetes
   kubectl exec -it <metrifuge-pod> -- nc -zv otel-collector 4317

   # Local
   nc -zv localhost 4317
   ```

2. Verify exporter configuration in Metrifuge CRDs:
   ```bash
   kubectl get exporters -o yaml
   ```

3. Check collector logs for errors:
   ```bash
   kubectl logs -l app=otel-collector --tail=100
   ```

### High Memory Usage

If the collector is consuming too much memory:

1. Adjust the memory limiter in `otel-collector-config.yaml`:
   ```yaml
   memory_limiter:
     limit_mib: 256  # Reduce from 512
     spike_limit_mib: 64  # Reduce from 128
   ```

2. Increase batch frequency to export more frequently:
   ```yaml
   batch:
     timeout: 5s  # Reduce from 10s
   ```

### No Data in Prometheus

1. Verify the Prometheus exporter endpoint:
   ```bash
   curl http://localhost:8889/metrics
   ```

2. Check that the metrics pipeline is configured correctly in `service.pipelines.metrics`

## Adding Additional Exporters

To forward data to external backends (e.g., Grafana Cloud, Datadog, New Relic):

1. Edit `otel-collector-config.yaml` and add the appropriate exporter:

```yaml
exporters:
  otlp/backend:
    endpoint: your-backend-endpoint:4317
    headers:
      api-key: ${BACKEND_API_KEY}
```

2. Add the exporter to the pipeline:

```yaml
service:
  pipelines:
    metrics:
      exporters: [logging, file/metrics, prometheus, otlp/backend]
```

3. Set environment variables in the deployment:

**Docker Compose:**
```yaml
environment:
  - BACKEND_API_KEY=your-api-key
```

**Kubernetes:**
```yaml
env:
  - name: BACKEND_API_KEY
    valueFrom:
      secretKeyRef:
        name: otel-secrets
        key: backend-api-key
```

## Production Considerations

For production deployments:

1. **Resource Limits**: Adjust CPU/memory limits based on load
2. **High Availability**: Run multiple collector replicas with a load balancer
3. **Persistent Storage**: Use persistent volumes for file exports
4. **Security**: Enable TLS for OTLP receivers
5. **Monitoring**: Set up alerting on collector health metrics
6. **Data Retention**: Configure appropriate log rotation policies

## References

- [OpenTelemetry Collector Documentation](https://opentelemetry.io/docs/collector/)
- [OTLP Specification](https://opentelemetry.io/docs/specs/otlp/)
- [Metrifuge Project README](../README.md)
