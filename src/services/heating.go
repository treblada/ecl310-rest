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
	"log"
	"math"
	"net/http"

	"github.com/treblada/ecl310-rest/generated/openapi"
	wrapper "github.com/treblada/ecl310-rest/modbus"
)

type HeatingApiService struct {
	openapi.HeatingApiService
	client wrapper.ZeroBasedAddressClientWrapper
}

type Int32Slice []int32

func (s Int32Slice) indexOf(x int32) int {
	for i := 0; i < len(s); i++ {
		if s[i] == x {
			return i
		}
	}
	return -1
}

func (s Int32Slice) has(x int32) bool {
	return s.indexOf(x) > -1
}

var validOutdoorTemps = Int32Slice{-30, -15, -5, 0, 5, 15}

func getSlopePnu(circuitNo int32) uint16 {
	return 10175 + uint16(circuitNo)*1000
}

func getMinMaxPnu(circuitNo int32) uint16 {
	return 10177 + uint16(circuitNo)*1000
}

func getTempCurvePointsPnu(circuitNo int32) uint16 {
	return 10400 + uint16(circuitNo)*1000
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

	slope := s.readPnu(getSlopePnu(circuitNo), 1)
	minMax := s.readPnu(getMinMaxPnu(circuitNo), 2)
	tempCurvePoints := s.readPnu(getTempCurvePointsPnu(circuitNo), 6)

	var curvePoints [6]openapi.FlowTempPoint
	for i := 0; i < len(validOutdoorTemps); i++ {
		curvePoints[i] = openapi.FlowTempPoint{
			OutdoorTemp: validOutdoorTemps[i],
			FlowTemp:    int32(binary.BigEndian.Uint16(tempCurvePoints[i*2 : i*2+2])),
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

func (s *HeatingApiService) SetHeatCurveBySlope(ctx context.Context, circuitNo int32, values openapi.SetHeatCurveBySlopeRequest) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	if circuitNo < 1 || circuitNo > 3 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid circuit number %d, not in [1,3]", circuitNo), nil))
	}

	if values.Slope > -0.1 || values.Slope < -10 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid slope value %f, must be in [-10, -0.1]", values.Slope), nil))
	}

	assertValidFlowTemperatureRange(values.MinFlowTemp, "min flow temp")
	assertValidFlowTemperatureRange(values.MaxFlowTemp, "max flow temp")

	if values.Slope != 0 {
		slopePnu := getSlopePnu(circuitNo)
		newSlopeInt := uint16(math.Round(float64(values.Slope) * -10))
		s.updateSinglePnu(slopePnu, newSlopeInt, "slope")
	}

	minMaxPnu := getMinMaxPnu(circuitNo)

	if values.MinFlowTemp != 0 {
		newMinTempInt := uint16(values.MinFlowTemp)
		s.updateSinglePnu(minMaxPnu, newMinTempInt, "min temp")
	}

	if values.MaxFlowTemp != 0 {
		newMaxTempInt := uint16(values.MaxFlowTemp)
		s.updateSinglePnu(minMaxPnu+1, newMaxTempInt, "max temp")
	}

	return s.GetHeatCurve(ctx, circuitNo)
}

func (s *HeatingApiService) SetHeatCurveByPoints(ctx context.Context, circuitNo int32, values openapi.SetHeatCurveByPointsRequest) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	if circuitNo < 1 || circuitNo > 3 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid circuit number %d, not in [1,3]", circuitNo), nil))
	}

	assertValidFlowTemperatureRange(values.MinFlowTemp, "min flow temp")
	assertValidFlowTemperatureRange(values.MaxFlowTemp, "max flow temp")

	for i := 0; i < len(values.CurvePoints); i++ {
		outTemp := values.CurvePoints[i].OutdoorTemp
		if !validOutdoorTemps.has(outTemp) {
			panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid outdoor temp %d, not in %v", outTemp, validOutdoorTemps), nil))
		}
		assertValidFlowTemperatureRange(values.CurvePoints[i].FlowTemp, fmt.Sprintf("flow temp for %d outside temp", outTemp))
	}

	minMaxPnu := getMinMaxPnu(circuitNo)

	if values.MinFlowTemp != 0 {
		newMinTempInt := uint16(values.MinFlowTemp)
		s.updateSinglePnu(minMaxPnu, newMinTempInt, "min temp")
	}

	if values.MaxFlowTemp != 0 {
		newMaxTempInt := uint16(values.MaxFlowTemp)
		s.updateSinglePnu(minMaxPnu+1, newMaxTempInt, "max temp")
	}

	tempCurvePointsPnu := getTempCurvePointsPnu(circuitNo)

	for _, curvePoint := range values.CurvePoints {
		i := validOutdoorTemps.indexOf(curvePoint.OutdoorTemp)
		s.updateSinglePnu(tempCurvePointsPnu+uint16(i), uint16(curvePoint.FlowTemp), fmt.Sprintf("%d outdoor temp", curvePoint.OutdoorTemp))
	}

	return s.GetHeatCurve(ctx, circuitNo)
}

func (s *HeatingApiService) readPnu(pnu uint16, quantity uint16) []byte {
	if result, err := s.client.ReadHoldingRegisters(pnu, quantity); err == nil {
		return result
	} else {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("Error reading PNU%d:%d", pnu, quantity), err))
	}
}

func (s *HeatingApiService) updateSinglePnu(pnu uint16, newValue uint16, label string) {
	oldValue := binary.BigEndian.Uint16(s.readPnu(pnu, 1))
	if oldValue != newValue {
		log.Printf("Updating %s: PNU%d:1 %d -> %d\n", label, pnu, oldValue, newValue)
		if _, err := s.client.WriteSingleRegister(pnu, newValue); err != nil {
			panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("Error writing %s PNU%d=%d", label, pnu, newValue), err))
		}
	}
}

func assertValidFlowTemperatureRange(tempValue int32, id string) {
	if tempValue != 0 && tempValue < 10 || tempValue > 150 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid value %d for %s. Valid values: [10, 150]", tempValue, id), nil))
	}
}
