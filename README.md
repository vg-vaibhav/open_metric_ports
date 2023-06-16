# Open Metrics Port Scanner

# Scans for open ports to be used for [Prometheus HTTP Service Discovery](https://prometheus.io/docs/prometheus/latest/http_sd/)

## How to run

```
go build
./open_metric_ports <host_address> <start_port_range> <end_port_range>
```

## HTTP Response

```
curl localhost:8888/hosts
```

```
{
  "targets": [
    "hostname.com:36036",
    "hostname.com:36018",
    "hostname.com:36006"
  ],
  "labels": { "key": "value" }
}
```

## TODO:

1. Add config file to start the server
2. Allow multiple servers to be scanned parallely
3. Allow multiple handlers params based on different host/group scanner
4. Add AWS api integration to autofetch list of hostnames
