package main

import (
	"fmt"

	"github.com/goburrow/modbus"
)

func main() {
	fmt.Println("Hello, World!")
	// Modbus TCP
	client := modbus.TCPClient("localhost:502")
	// Read input register 9
	results, err := client.ReadHoldingRegisters(278, 4)
	if err == nil {
		fmt.Println(results)
	} else {
		fmt.Println(err)
	}
}
