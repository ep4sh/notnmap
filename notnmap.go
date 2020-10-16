package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	MAXPORTS  = 50000 // 1-65535
	WORKERS   = 100   // 5-100
	TIMEOUT   = 100   // msec
	RESULTBUF = 512
)

func main() {
	fmt.Println(`
	 _____  _____  _____  _____  _____  _____  _____
	 |   | ||     ||_   _||   | ||     ||  _  ||  _  |
	 | | | ||  |  |  | |  | | | || | | ||     ||   __|
	 |_|___||_____|  |_|  |_|___||_|_|_||__|__||__|

		`)

	ports := make(chan int, MAXPORTS)
	results := make(chan int, RESULTBUF)
	done := make(chan bool)
	var wg sync.WaitGroup

	// Add ports to channel
	for i := 1; i <= MAXPORTS; i++ {
		ports <- i
	}
	close(ports)

	// Run workers
	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go func(i int) {
			wg.Done()
			worker(i, ports, results, done)
		}(i)
	}
	wg.Wait()
	for {
		select {
		case openPort, ok := <-results:
			fmt.Println("Found open port:", openPort)
			if !ok {
				// ROFL
				fmt.Println("Err, please restart pc")
			}
		case <-done:
			os.Exit(0)
		}
	}
}

func worker(id int, ports <-chan int, results chan<- int, done chan<- bool) {
	for {
		port, more := <-ports
		if more {
			host := fmt.Sprintf("%s:%d", os.Args[1], port)
			conn, err := net.DialTimeout("tcp", host, TIMEOUT*time.Millisecond)
			if err != nil {
				continue
			}
			results <- port
			conn.Close()
		} else {
			done <- true
			return
		}
	}
}
