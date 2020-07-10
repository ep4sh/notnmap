package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	MAXPORT = 65536
)

func worker(ports, results chan int) {
	for port := range ports {
		host := fmt.Sprintf("%s:%d", os.Args[1], port)
		conn, err := net.DialTimeout("tcp", host, 10*time.Second)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- port
		return
	}
}

func main() {
	if len(os.Args) >= 2 {
		ports := make(chan int, MAXPORT)
		results := make(chan int)

		for i := 0; i <= 1000; i++ {
			go worker(ports, results)
		}
		go func() {
			for i := 0; i <= cap(ports); i++ {
				ports <- i
			}
		}()
		for i := 1; i <= MAXPORT; i++ {
			port := <-results
			if port != 0 {
				fmt.Printf("Found open port: %d \n", port)
			}
		}
		close(ports)
		close(results)
	} else {
		log.Fatalf("Specify host to scan")
	}
}
