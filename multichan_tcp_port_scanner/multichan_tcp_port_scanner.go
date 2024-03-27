package main

import (
	"fmt"
	"net"
	"sort"
)

var host string
var network string
var ports_rage int

func worker(ports, results chan int) {
	// тут инициализируем логику обработки каннала при будущем заполнении через другую горутину
	for port := range ports {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err_net_dial := net.Dial(network, address)

		// в случае если мы не сможем получитьт ответ с порта, то мы заполняем канал результатов 0, чтобы не оставлять висеть канал не обработанным
		if err_net_dial != nil {
			results <- 0
			continue
		}

		conn.Close()
		results <- port
	}
}

func main() {
	fmt.Println("Hello! It's TCP-client for port scrinning.")

	fmt.Println("Please, enter network (default tcp): ")
	_, err_network := fmt.Scanln(&network)
	if err_network != nil {
		network = "tcp"
	}

	fmt.Printf("Please, enter host (default localhost): ")
	_, err_host := fmt.Scanln(&host)
	if err_host != nil {
		host = "localhost"
	}

	fmt.Printf("Please, enter how many posts scan do you need (default 65535): ")
	_, err_ports_rage := fmt.Scanln(&ports_rage)
	if err_ports_rage != nil {
		ports_rage = 65535
	}

	ports := make(chan int, 100)
	results := make(chan int)

	// создаем горутины инициализируя 100 воркеров, где каждый каннал вмещает в себя 100 элементов
	for channel := 0; channel < cap(ports); channel++ {
		go worker(ports, results)
	}

	fmt.Printf("Scan for %s start...\n", host)
	go func() {
		for port := 1; port <= ports_rage; port++ {
			// здесь мы заполняем канал числами по которым будет итерироваться воркер где будет использоваться канал
			ports <- port

			fmt.Printf("try connection port: %d", port)
			fmt.Printf("\r")
		}
	}()

	var openports []int
	var success_counter = 0

	for i := 0; i < ports_rage; i++ {
		port := <-results
		if port != 0 {
			success_counter++
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)

	fmt.Printf("\n")
	if success_counter == 0 {
		fmt.Println("No open ports")
	} else {
		for _, port := range openports {
			fmt.Printf("%d open\n", port)
		}
	}
}
