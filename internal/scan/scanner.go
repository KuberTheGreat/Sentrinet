package scan

import (
	"fmt"
	"net"
	"time"
)

type PortResult struct{
	Port int
	IsOpen bool
	Duration int64
}

func ScanPort(target string, port int) PortResult{
	start := time.Now()
	address := net.JoinHostPort(target, fmt.Sprintf("%d", port))

	addrs, err := net.LookupHost(target)
	if err != nil {
		fmt.Printf("[debug] LookupHost %s error: %v\n", target, err)
	} else {
		fmt.Printf("[debug] %s resolves to %v\n", target, addrs)
	}
	
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	duration := time.Since(start).Milliseconds()

	if err != nil{
		return PortResult{Port: port, IsOpen: false, Duration: duration}
	}

	conn.Close()
	return PortResult{Port: port, IsOpen: true, Duration: duration}
}

func ScanRange(target string, startPort, endPort int) []PortResult{
	results := make([]PortResult, 0)
	ch := make(chan PortResult)

	for port := startPort; port <= endPort; port++{
		go func(p int){
			ch <- ScanPort(target, p)
		}(port)
	}

	for i := startPort; i <= endPort; i++{
		result := <-ch
		results = append(results, result)
	}

	return results
}