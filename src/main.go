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
	"github.com/treblada/ecl310-rest/generated/openapi"

	"log"
	"net/http"

	wrapper "github.com/treblada/ecl310-rest/modbus"
	api "github.com/treblada/ecl310-rest/services"
)

func main() {
	fmt.Println("ECL310 API starting")
	config := parseCmdLine()
	// Modbus TCP
	modbusClient := wrapper.NewZeroBasedAddressClientWrapper(modbus.TCPClient(fmt.Sprintf("%s:%d", config.eclHost, config.eclPort)))
	log.Printf("ECL client ready.")

	HealthService := api.NewHealthApiService(&modbusClient)
	HealthServiceController := openapi.NewHealthApiController(HealthService)

	SystemService := api.NewSystemApiService(&modbusClient)
	SystemServiceController := openapi.NewSystemApiController(SystemService)

	router := openapi.NewRouter(HealthServiceController, SystemServiceController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
