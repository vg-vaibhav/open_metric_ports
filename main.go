package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrInvalidPort = errors.New("")

func main() {
	hostname := os.Args[1]
	portRangeStart, _ := strconv.Atoi(os.Args[2])
	portRangeEnd, _ := strconv.Atoi(os.Args[3])
	response := Response{
		Targets: []string{},
		Labels:  map[string]string{"key": "value"},
	}

	server := Server{clientAddress: "localhost:8888"}
	var wg sync.WaitGroup

	wg.Add(1)
	go server.Apilisten(&response)

	wg.Add(1)
	response.scan(hostname, portRangeStart, portRangeEnd)

	wg.Wait()
}

func (r *Response) scan(hostname string, portStart, portEnd int) {
	go func(hostname string) {
		for {
			select {
			case <-time.After(10 * time.Second):
				log.Info("Scannning")
				ports, err := getPortsInUse(hostname)
				if err != nil {
					log.Errorf("Failed while scanning ports %v", err)
					continue
				}
				validPorts := []int{}
				for _, v := range ports {
					if v >= portStart && v <= portEnd {
						validPorts = append(validPorts, v)
					}
				}
				metricPorts := scanForMetricsHandler(hostname, validPorts)
				hosts := []string{}
				for _, port := range metricPorts {
					host := fmt.Sprintf("%s:%d", hostname, port)
					hosts = append(hosts, host)
				}
				r.Targets = hosts
			}
		}
	}(hostname)
}

func scanForMetricsHandler(ip string, ports []int) []int {
	// HTTP client with timeout of 10 seconds
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	metricPorts := []int{}
	queue := make(chan int, 1)

	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		// Send HTTP request in a separate goroutine
		go func(ip string, port int) {
			defer wg.Done()

			resp, err := client.Get(fmt.Sprintf("http://%s:%d/metrics", ip, port))
			if err != nil {
				// fmt.Printf("Error connecting to %s:%d: %s\n", vm, port, err)
				return
			}

			if resp.StatusCode == 200 {
				_, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response from %s:%d: %s\n", ip, port, err)
					return
				}
				queue <- port
			}
		}(ip, port)
	}

	go func() {
		for port := range queue {
			metricPorts = append(metricPorts, port)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(queue)
	<-queue
	return metricPorts
}

func getPortsInUse(hostname string) ([]int, error) {
	cmd := exec.Command("ssh", hostname, "netstat", "-lnp", "--tcp") // Execute netstat command with appropriate flags
	output, err := cmd.Output()
	if err != nil {
		return []int{}, err
	}

	ports := parseBoundPorts(string(output))
	return ports, nil

}

func parseBoundPorts(output string) []int {
	lines := strings.Split(output, "\n")
	ports := make([]int, 0)
	for _, line := range lines {
		fields := strings.Fields(line)

		port, err := parseTCP(fields)
		if err == nil {
			ports = append(ports, port)
		}

		port, err = parseTCP6(fields)
		if err == nil {
			ports = append(ports, port)
		}

	}
	return ports
}

func parseTCP6(fields []string) (int, error) {
	if len(fields) >= 4 && fields[0] == "tcp6" && strings.HasPrefix(fields[5], "LISTEN") {
		address := fields[3]
		return strconv.Atoi(strings.Split(address, ":")[3])
	}
	return 0, ErrInvalidPort
}

func parseTCP(fields []string) (int, error) {
	if len(fields) >= 4 && fields[0] == "tcp" && strings.HasPrefix(fields[5], "LISTEN") {
		address := fields[3]
		return strconv.Atoi(strings.Split(address, ":")[1])
	}
	return 0, ErrInvalidPort
}
