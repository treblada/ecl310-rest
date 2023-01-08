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
	log.Println("ECL310 API starting")
	config := parseCmdLine()
	log.Printf("Working with remote instance %s:%d\n", config.eclHost, config.eclPort)
	// Modbus TCP
	modbusClient := wrapper.NewModbusClientWrapper(modbus.TCPClient(fmt.Sprintf("%s:%d", config.eclHost, config.eclPort)))
	log.Println("ECL client ready.")

	HealthService := api.NewHealthApiService(&modbusClient)
	HealthServiceController := openapi.NewHealthApiControllerWithErrorHandler(HealthService, api.ApiErrorHandler)

	SystemService := api.NewSystemApiService(&modbusClient)
	SystemServiceController := openapi.NewSystemApiControllerWithErrorHandler(SystemService, api.ApiErrorHandler)

	HeatingService := api.NewHeatingApiService(&modbusClient)
	HeatingServiceController := openapi.NewHeatingApiControllerWithErrorHandler(HeatingService, api.ApiErrorHandler)

	router := openapi.NewRouter(HealthServiceController, SystemServiceController, HeatingServiceController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
