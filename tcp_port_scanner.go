package main

import (
	"fmt"
	"net"
	"sync"
)

var host string
var network string

func worker(ports chan int, wg *sync.WaitGroup) {
	for port := range ports {
		address := fmt.Sprintf("%s:%d", host, port)

		conn, err_net_dial := net.Dial(network, address)

		if err_net_dial != nil {
			continue
		}

		conn.Close()
		fmt.Printf("Connection success for port: %d\n", port)
		wg.Done()
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("Hello! It's TCP-client for port scrinning.")

	fmt.Println("Please, enter network (default tcp): ")
	_, err_network := fmt.Scanln(&network)
	if err_network != nil {
		network = "tcp"
	}

	fmt.Println("Please, enter host (default scanme.nmap.org): ")
	_, err_host := fmt.Scanln(&host)
	if err_host != nil {
		host = "scanme.nmap.org"
	}

	fmt.Printf("Scan for %s start...\n", host)

	ports := make(chan int, 100)

	for channel := 0; channel < cap(ports); channel++ {
		go worker(ports, &wg)
	}

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		ports <- port
	}

	wg.Wait()
	close(ports)
}
