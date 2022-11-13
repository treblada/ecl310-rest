/*
This file is part of ecl310-rest.

ecl310-rest is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

ecl310-rest is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
ecl310-rest. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"

	"github.com/goburrow/modbus"
	openapi "github.com/treblada/ecl310-rest/openapi/go"

	"log"
	"net/http"
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

	log.Printf("Server started")

	HealthService := openapi.NewHealthApiService()
	HealthServiceController := openapi.NewHealthApiController(HealthService)

	router := openapi.NewRouter(HealthServiceController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
