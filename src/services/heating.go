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
	"context"
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/treblada/ecl310-rest/generated/openapi"
	wrapper "github.com/treblada/ecl310-rest/modbus"
)

type HeatingApiService struct {
	openapi.HeatingApiService
	client wrapper.ZeroBasedAddressClientWrapper
}

func NewHeatingApiService(client wrapper.ZeroBasedAddressClientWrapper) openapi.HeatingApiServicer {
	if client == nil {
		panic("No modbus client provided for System API service")
	}
	return &HeatingApiService{
		client: client,
	}
}

func (s *HeatingApiService) GetHeatCurve(ctx context.Context, circuitNo int32) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	if circuitNo < 1 || circuitNo > 3 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid circuit number %d, not in [1,3]", circuitNo), nil))
	}

	var slope []byte
	var minMax []byte
	var temperature []byte
	var slopePnu uint16 = 10175 + uint16(circuitNo)*1000
	var minMaxPnu uint16 = 10177 + uint16(circuitNo)*1000
	var temperaturePnu uint16 = 10400 + uint16(circuitNo)*1000

	var err error

	if slope, err = s.client.ReadHoldingRegisters(slopePnu, 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d", slopePnu), err))
	}
	if minMax, err = s.client.ReadHoldingRegisters(minMaxPnu, 2); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d+2", minMaxPnu), err))
	}
	if temperature, err = s.client.ReadHoldingRegisters(temperaturePnu, 6); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d+6", temperaturePnu), err))
	}

	var outdoorTemps = []int32{-30, -15, -5, 0, 5, 15}
	var curvePoints [6]openapi.FlowTempPoint
	for i := 0; i < 6; i++ {
		curvePoints[i] = openapi.FlowTempPoint{
			OutdoorTemp: outdoorTemps[i],
			FlowTemp:    int32(binary.BigEndian.Uint16(temperature[i*2 : i*2+2])),
		}
	}

	body := openapi.GetHeatCurveResponse{
		Slope:       float32(binary.BigEndian.Uint16(slope)) / -10.0,
		MinFlowTemp: int32(binary.BigEndian.Uint16(minMax[0:2])),
		MaxFlowTemp: int32(binary.BigEndian.Uint16(minMax[2:4])),
		CurvePoints: curvePoints[:],
	}

	return openapi.Response(200, body), nil
}
