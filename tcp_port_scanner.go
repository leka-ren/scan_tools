package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var host string
	var network string

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

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		go func(port_p int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", host, port_p)

			conn, err_net_dial := net.Dial(network, address)

			if err_net_dial != nil {
				return
			}

			conn.Close()
			fmt.Printf("Connection success for port: %d\n", port_p)
		}(port)
	}
	wg.Wait()
}
