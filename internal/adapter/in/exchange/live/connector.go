package live

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Quote struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func ConnectAndRead(port int, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := "exchange1:40101"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Error connecting to port %d: %v", port, err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to port %d", port)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		var q Quote
		err := json.Unmarshal([]byte(line), &q)
		if err != nil {
			log.Printf("Error parsing JSON on port %d: %v", port, err)
			continue
		}
		fmt.Printf("Port %d: %+v\n", port, q)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from port %d: %v", port, err)
	}
}
