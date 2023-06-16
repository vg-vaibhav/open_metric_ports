# Open Metrics Port Scanner

# Scans for open ports for Prometheus/Micrometer metrics

## How to run

```
go build
./open_metric_ports <host_address> <start_port_range> <end_port_range>
```

## TODO:

1. Add config file to start the server
2. Allow multiple servers to be scanned parallely
3. Allow multiple handlers params based on different host/group scanner
4. Add AWS api integration to autofetch list of hostnames
