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

package api

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/goburrow/modbus"
	"github.com/treblada/ecl310-rest/generated/openapi"
)

func ApiErrorHandler(w http.ResponseWriter, r *http.Request, err error, result *openapi.ImplResponse) {
	if _, ok := err.(*net.OpError); ok {
		// Network connection error
		openapi.EncodeJSONResponse(err.Error(), func(i int) *int { return &i }(http.StatusBadGateway), w)
	} else if _, ok := err.(*modbus.ModbusError); ok {
		// Modbus communication error
		openapi.EncodeJSONResponse(err.Error(), func(i int) *int { return &i }(http.StatusBadGateway), w)
	} else if typedErr, ok := err.(*ApiError); ok {
		log.Println(err.Error())
		openapi.EncodeJSONResponse(fmt.Sprintf("%s; %s", typedErr.Message, typedErr.Cause.Error()), &typedErr.Code, w)
	} else {
		openapi.DefaultErrorHandler(w, r, err, result)
	}
}
