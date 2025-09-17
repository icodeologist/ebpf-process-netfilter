package main

import (
	"fmt"
	"net"
	"time"
)

func testConnection(port string) {
	fmt.Printf("Testing port %s...\n", port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", "8.8.8.8:"+port, 3*time.Second)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Port %s FAILED after %v: %v\n", port, duration, err)
	} else {
		fmt.Printf("Port %s CONNECTED after %v\n", port, duration)
		conn.Close()
	}
	fmt.Println("---")
}

func main() {
	testConnection("4040") // Should be allowed (but timeout)
	testConnection("53")   // Should be blocked (DNS port, usually fast)
	testConnection("80")   // Should be blocked
}
