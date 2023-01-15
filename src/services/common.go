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
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"github.com/treblada/ecl310-rest/generated/openapi"
	wrapper "github.com/treblada/ecl310-rest/modbus"
)

func handlePanic(panic any) (response openapi.ImplResponse, funcErr error) {
	response = openapi.ImplResponse{}
	if typedPanic, ok := panic.(error); ok {
		funcErr = typedPanic
	} else {
		funcErr = fmt.Errorf("%v", panic)
	}
	return
}

func readPnu(c wrapper.ZeroBasedAddressClientWrapper, pnu uint16, quantity uint16) []byte {
	if result, err := c.ReadHoldingRegisters(pnu, quantity); err == nil {
		return result
	} else {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("Error reading PNU%d:%d", pnu, quantity), err))
	}
}

func updateSinglePnu(c wrapper.ZeroBasedAddressClientWrapper, pnu uint16, newValue uint16, label string) {
	oldValue := binary.BigEndian.Uint16(readPnu(c, pnu, 1))
	if oldValue != newValue {
		log.Printf("Updating %s: PNU%d:1 %d -> %d\n", label, pnu, oldValue, newValue)
		if _, err := c.WriteSingleRegister(pnu, newValue); err != nil {
			panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("Error writing %s PNU%d=%d", label, pnu, newValue), err))
		}
	}
}
