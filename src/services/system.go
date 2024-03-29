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
	"net/http"

	"github.com/treblada/ecl310-rest/generated/openapi"
	wrapper "github.com/treblada/ecl310-rest/modbus"
)

type SystemApiService struct {
	openapi.SystemApiService
	client wrapper.ZeroBasedAddressClientWrapper
}

func NewSystemApiService(client wrapper.ZeroBasedAddressClientWrapper) openapi.SystemApiServicer {
	if client == nil {
		panic("No modbus client provided for System API service")
	}
	return &SystemApiService{
		client: client,
	}
}

func (s *SystemApiService) GetSystemInfo(ctx context.Context) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	var pnu19 []byte
	var pnu34_37 []byte
	var pnu258 []byte
	var pnu278_289 []byte
	var pnu2060_2063 []byte
	var pnu2099 []byte
	var err error

	if pnu19, err = s.client.ReadHoldingRegisters(19, 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU19", err))
	}
	if pnu34_37, err = s.client.ReadHoldingRegisters(34, 4); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU34+4", err))
	}
	if pnu258, err = s.client.ReadHoldingRegisters(258, 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU258", err))
	}
	if pnu278_289, err = s.client.ReadHoldingRegisters(278, 12); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU278+12", err))
	}
	if pnu2060_2063, err = s.client.ReadHoldingRegisters(2060, 4); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU2060+4", err))
	}
	if pnu2099, err = s.client.ReadHoldingRegisters(2099, 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU2099", err))
	}

	body := openapi.GetSystemInfoResponse{
		HardwareRevision: fmt.Sprintf("087H%d", binary.BigEndian.Uint16(pnu19)),
		SoftwareVersion:  int32(binary.BigEndian.Uint16(pnu34_37[2:4])),
		SerialNumber:     int64(binary.BigEndian.Uint32(pnu34_37[4:8])),
		AddressType:      decodeAddressType(pnu258),
		IpAddress: fmt.Sprintf(
			"%d.%d.%d.%d",
			binary.BigEndian.Uint16(pnu278_289[0:2]),
			binary.BigEndian.Uint16(pnu278_289[2:4]),
			binary.BigEndian.Uint16(pnu278_289[4:6]),
			binary.BigEndian.Uint16(pnu278_289[6:8]),
		),
		Netmask: fmt.Sprintf(
			"%d.%d.%d.%d",
			binary.BigEndian.Uint16(pnu278_289[16:18]),
			binary.BigEndian.Uint16(pnu278_289[18:20]),
			binary.BigEndian.Uint16(pnu278_289[20:22]),
			binary.BigEndian.Uint16(pnu278_289[22:24]),
		),
		Gateway: fmt.Sprintf(
			"%d.%d.%d.%d",
			binary.BigEndian.Uint16(pnu278_289[8:10]),
			binary.BigEndian.Uint16(pnu278_289[10:12]),
			binary.BigEndian.Uint16(pnu278_289[12:14]),
			binary.BigEndian.Uint16(pnu278_289[14:16]),
		),
		Application:        decodeApplicationName(pnu2060_2063),
		ApplicationVersion: fmt.Sprintf("%d.%d", pnu2060_2063[6], pnu2060_2063[7]),
		ProductionYear:     2000 + int32(pnu2099[0]),
		ProductionWeek:     int32(pnu2099[1]),
	}
	return openapi.Response(http.StatusOK, body), nil
}

func decodeAddressType(pnu258 []byte) string {
	switch dhcpFlag := binary.BigEndian.Uint16(pnu258); dhcpFlag {
	case 0:
		return "DHCP"
	case 1:
		return "STATIC"
	default:
		panic(fmt.Errorf("invalid address type %d on PNU 258", dhcpFlag))
	}
}

func decodeApplicationName(pnu2060_2063 []byte) string {
	appPrefix := string(rune(binary.BigEndian.Uint16(pnu2060_2063[0:2])))
	appType := binary.BigEndian.Uint16(pnu2060_2063[2:4])
	appSubType := binary.BigEndian.Uint16(pnu2060_2063[4:6])
	return fmt.Sprintf("%v%d.%d", appPrefix, appType, appSubType)
}

func (s *SystemApiService) GetSystemCircuit(ctx context.Context, circuitNo int32) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	if circuitNo < 1 || circuitNo > 3 {
		return openapi.ImplResponse{}, NewApiError(400, "Circuit number must be in [1-3]", nil)
	}

	var circMode []byte
	var circState []byte
	var err error

	modeAddr := 4200 + circuitNo
	stateAddr := 4210 + circuitNo

	if circMode, err = s.client.ReadHoldingRegisters(uint16(modeAddr), 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d", modeAddr), err))
	}
	if circState, err = s.client.ReadHoldingRegisters(uint16(stateAddr), 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d", stateAddr), err))
	}

	body := openapi.GetSystemCircuitResponse{
		Mode:   GetCircuitMode(binary.BigEndian.Uint16(circMode)).String(),
		Status: GetCircuitState(binary.BigEndian.Uint16(circState)).String(),
	}
	return openapi.Response(200, body), nil
}

func (s *SystemApiService) GetSystemCircuits(ctx context.Context) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	var circModes []byte
	var circStates []byte
	var err error

	var modeBaseAddr uint16 = 4200
	var stateBaseAddr uint16 = 4210

	if circModes, err = s.client.ReadHoldingRegisters(modeBaseAddr+1, 2); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU4201+2", err))
	}
	if circStates, err = s.client.ReadHoldingRegisters(stateBaseAddr+1, 2); err != nil {
		panic(NewApiError(http.StatusBadGateway, "PNU4211+2", err))
	}

	heating := openapi.GetSystemCircuitResponse{
		Mode:   GetCircuitMode(binary.BigEndian.Uint16(circModes[:2])).String(),
		Status: GetCircuitState(binary.BigEndian.Uint16(circStates[:2])).String(),
	}
	warmWater := openapi.GetSystemCircuitResponse{
		Mode:   GetCircuitMode(binary.BigEndian.Uint16(circModes[2:4])).String(),
		Status: GetCircuitState(binary.BigEndian.Uint16(circStates[2:4])).String(),
	}

	// this should be a pointer, unfortunatelly the openapi-generator does not seem to support it
	circ3 := openapi.GetSystemCircuitResponse{}

	if circModes, err = s.client.ReadHoldingRegisters(modeBaseAddr+3, 1); err == nil {
		if circStates, err = s.client.ReadHoldingRegisters(stateBaseAddr+3, 1); err == nil {
			circ3 = openapi.GetSystemCircuitResponse{
				Mode:   GetCircuitMode(binary.BigEndian.Uint16(circModes[:2])).String(),
				Status: GetCircuitState(binary.BigEndian.Uint16(circStates[:2])).String(),
			}
		} else {
			log.Printf("Error reading PNU 4213: %v", err)
		}
	} else {
		log.Printf("Error reading PNU 4203: %v", err)
	}

	body := openapi.GetSystemCircuitsResponse{
		Heating:   heating,
		WarmWater: warmWater,
		Circuit3:  circ3,
	}

	return openapi.Response(200, body), nil
}

func (s *SystemApiService) GetSystemDateTime(ctx context.Context) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	body := s.getDateTime()
	return openapi.Response(http.StatusOK, body), nil
}

var pnuHour uint16 = 64045
var pnuMinute uint16 = 64046
var pnuDay uint16 = 64047
var pnuMonth uint16 = 64048
var pnuYear uint16 = 64049
var pnuDst uint16 = 10198

func (s *SystemApiService) getDateTime() openapi.GetSystemDateTime {
	var datetime []byte
	var dst []byte
	var err error

	if datetime, err = s.client.ReadHoldingRegisters(pnuHour, 5); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d:5", pnuHour), err))
	}

	if dst, err = s.client.ReadHoldingRegisters(pnuDst, 1); err != nil {
		panic(NewApiError(http.StatusBadGateway, fmt.Sprintf("PNU%d:1", pnuDst), err))
	}

	return openapi.GetSystemDateTime{
		Hour:               int32(binary.BigEndian.Uint16(datetime[0:2])),
		Minute:             int32(binary.BigEndian.Uint16(datetime[2:4])),
		Day:                int32(binary.BigEndian.Uint16(datetime[4:6])),
		Month:              int32(binary.BigEndian.Uint16(datetime[6:8])),
		Year:               int32(binary.BigEndian.Uint16(datetime[8:10])),
		AutoDaylightSaving: binary.BigEndian.Uint16(dst) == uint16(1),
	}
}

func (s *SystemApiService) SetSystemDateTime(ctx context.Context, newDateTime openapi.GetSystemDateTime) (response openapi.ImplResponse, funcErr error) {
	defer func() {
		if panic := recover(); panic != nil {
			response, funcErr = handlePanic(panic)
		}
	}()

	if newDateTime.Year < 2009 || newDateTime.Year > 2099 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid year %d [2009, 2099]", newDateTime.Year), nil))
	}
	if newDateTime.Month < 1 || newDateTime.Month > 12 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid month %d [1, 12]", newDateTime.Month), nil))
	}
	if newDateTime.Hour < 0 || newDateTime.Hour > 23 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid hour %d [0, 23]", newDateTime.Hour), nil))
	}
	if newDateTime.Minute < 0 || newDateTime.Minute > 59 {
		panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid minute %d [0, 59]", newDateTime.Minute), nil))
	}

	daysPerMonth := map[int32]int32{1: 31, 2: 29, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}

	if daysPerMonth[newDateTime.Month] < newDateTime.Day {
		panic(NewApiError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid day %d for month %d [%d]", newDateTime.Day, newDateTime.Month, daysPerMonth[newDateTime.Month]),
			nil,
		))
	}

	if newDateTime.Month == 2 {
		if newDateTime.Day > 28 && !isLeapYear(newDateTime.Year) {
			panic(NewApiError(http.StatusBadRequest, fmt.Sprintf("Invalid day %d for month %d", newDateTime.Day, newDateTime.Month), nil))
		}
	}

	now := s.getDateTime()

	if newDateTime.Month == 2 && newDateTime.Day == 29 {
		// must be a leap year, otherwise we would have triggered a panic before
		// year, day, month
		updateSinglePnu(s.client, pnuYear, uint16(newDateTime.Year), "year")
		updateSinglePnu(s.client, pnuDay, uint16(newDateTime.Day), "day")
		updateSinglePnu(s.client, pnuMonth, uint16(newDateTime.Month), "month")
	} else if daysPerMonth[newDateTime.Month] > daysPerMonth[now.Month] {
		// month, day, year
		updateSinglePnu(s.client, pnuMonth, uint16(newDateTime.Month), "month")
		updateSinglePnu(s.client, pnuDay, uint16(newDateTime.Day), "day")
		updateSinglePnu(s.client, pnuYear, uint16(newDateTime.Year), "year")
	} else {
		// day, month, year
		updateSinglePnu(s.client, pnuDay, uint16(newDateTime.Day), "day")
		updateSinglePnu(s.client, pnuMonth, uint16(newDateTime.Month), "month")
		updateSinglePnu(s.client, pnuYear, uint16(newDateTime.Year), "year")
	}

	updateSinglePnu(s.client, pnuHour, uint16(newDateTime.Hour), "hour")
	updateSinglePnu(s.client, pnuMinute, uint16(newDateTime.Minute), "minute")

	updateSinglePnu(s.client, pnuDst, boolToUint16(newDateTime.AutoDaylightSaving), "DST")

	return s.GetSystemDateTime(ctx)
}

func isLeapYear(year int32) bool {
	// we ignore 100/400 year rules, b/c valid range is 2009-2099
	return year%4 == 0
}

func boolToUint16(value bool) uint16 {
	if value {
		return 1
	} else {
		return 0
	}
}
