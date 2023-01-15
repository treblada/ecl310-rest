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

package api_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/treblada/ecl310-rest/generated/openapi"
	"github.com/treblada/ecl310-rest/mocks"
	api "github.com/treblada/ecl310-rest/services"
	"gotest.tools/v3/assert"
)

func TestGetSystemInfo__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch address {
			case 19:
				return []byte{0, 42}, nil
			case 34:
				return []byte{
					0, 0,
					1, 3,
					0, 2, 0, 42,
				}, nil
			case 258:
				return []byte{0, 1}, nil
			case 278:
				return []byte{
					0, 192, 0, 168, 0, 1, 0, 1,
					0, 10, 0, 0, 0, 0, 0, 1,
					0, 255, 0, 255, 0, 255, 0, 128,
				}, nil
			case 2060:
				return []byte{
					0, 'f',
					0, 42,
					1, 1,
					3, 7,
				}, nil
			case 2099:
				return []byte{21, 33}, nil
			default:
				t.Logf("Failed access to %d:%d", address, quantity)
				t.Fail()
				return nil, errors.New("Mock failure")
			}
		},
	}
	service := api.NewSystemApiService(mock)
	response, err := service.GetSystemInfo(context.TODO())
	assert.NilError(t, err)
	assert.Equal(t, 200, response.Code)
	body := response.Body.(openapi.GetSystemInfoResponse)
	assert.Equal(t, "087H42", body.HardwareRevision)
	assert.Equal(t, int32(1<<8+3), body.SoftwareVersion)
	assert.Equal(t, int64(2<<16+42), body.SerialNumber)
	assert.Equal(t, "STATIC", body.AddressType)
	assert.Equal(t, "192.168.1.1", body.IpAddress)
	assert.Equal(t, "255.255.255.128", body.Netmask)
	assert.Equal(t, "10.0.0.1", body.Gateway)
	assert.Equal(t, "f42.257", body.Application)
	assert.Equal(t, "3.7", body.ApplicationVersion)
	assert.Equal(t, int32(2021), body.ProductionYear)
	assert.Equal(t, int32(33), body.ProductionWeek)
}

func TestGetSystemInfo__failure(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			return nil, errors.New("Mock error")
		},
	}
	service := api.NewSystemApiService(mock)
	_, err := service.GetSystemInfo(context.TODO())
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == 502)
}

func TestGetSystemCircuit__failure(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			return nil, errors.New("Mock error")
		},
	}
	service := api.NewSystemApiService(mock)
	_, err := service.GetSystemCircuit(context.TODO(), 1)
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadGateway)
}

func TestGetSystemCircuit__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			assert.Check(t, quantity == 1)
			switch address {
			case 4201:
				return []byte{0, 0}, nil
			case 4211:
				return []byte{0, 1}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
	}
	service := api.NewSystemApiService(mock)
	response, err := service.GetSystemCircuit(context.TODO(), 1)
	assert.NilError(t, err)
	assert.Check(t, http.StatusOK == response.Code)
	body := response.Body.(openapi.GetSystemCircuitResponse)
	assert.Check(t, body.Mode == api.Manual.String())
	assert.Check(t, body.Status == api.PreComfort.String())
}

func TestGetSystemCircuit__invalidRequestParam(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			t.Error("Should not have been called")
			t.FailNow()
			return nil, errors.New("Failure")
		},
	}
	service := api.NewSystemApiService(mock)
	var err error
	var apiErr *api.ApiError

	_, err = service.GetSystemCircuit(context.TODO(), 0)
	apiErr = err.(*api.ApiError)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)

	_, err = service.GetSystemCircuit(context.TODO(), 4)
	apiErr = err.(*api.ApiError)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestGetSystemCircuits__failure(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			return nil, errors.New("Mock error")
		},
	}
	service := api.NewSystemApiService(mock)
	_, err := service.GetSystemCircuits(context.TODO())
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadGateway)
}

func TestGetSystemCircuits__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch address {
			case 4201:
				assert.Equal(t, uint16(2), quantity)
				return []byte{0, 2, 0, 3}, nil
			case 4211:
				assert.Equal(t, uint16(2), quantity)
				return []byte{0, 1, 0, 2}, nil
			case 4203:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 4}, nil
			case 4213:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 3}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
	}
	service := api.NewSystemApiService(mock)
	response, err := service.GetSystemCircuits(context.TODO())
	assert.NilError(t, err)
	assert.Check(t, response.Code == http.StatusOK)
	body := response.Body.(openapi.GetSystemCircuitsResponse)
	assert.Check(t, body.Heating.Mode == api.ConstantComfortTemp.String())
	assert.Check(t, body.Heating.Status == api.PreComfort.String())
	assert.Check(t, body.WarmWater.Mode == api.ConstantSetbackTemp.String())
	assert.Check(t, body.WarmWater.Status == api.Comfort.String())
	assert.Check(t, body.Circuit3.Mode == api.FrostProtection.String())
	assert.Check(t, body.Circuit3.Status == api.PreSetback.String())
}

func TestGetSystemDateTime__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch address {
			case 64045:
				assert.Equal(t, uint16(5), quantity)
				return []byte{0, 10, 0, 11, 0, 14, 0, 2, 7, 230}, nil
			case 10198:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 1}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
	}
	service := api.NewSystemApiService(mock)
	response, err := service.GetSystemDateTime(context.TODO())
	assert.NilError(t, err)
	assert.Check(t, response.Code == http.StatusOK)
	body := response.Body.(openapi.GetSystemDateTime)
	assert.Equal(t, body.Hour, int32(10))
	assert.Equal(t, body.Minute, int32(11))
	assert.Equal(t, body.Day, int32(14))
	assert.Equal(t, body.Month, int32(2))
	assert.Equal(t, body.Year, int32(2022))
	assert.Check(t, body.AutoDaylightSaving)
}

func TestSetSystemDateTime__invalidYear(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 1999, Month: 2, Day: 1, Hour: 12, Minute: 2, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "year 1999")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__invalidMonth(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2009, Month: 13, Day: 1, Hour: 12, Minute: 2, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "month 13")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__invalidDay(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2009, Month: 2, Day: 32, Hour: 12, Minute: 2, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "day 32")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__invalidHour(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2009, Month: 2, Day: 1, Hour: 24, Minute: 2, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "hour 24")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__invalidMinute(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2009, Month: 2, Day: 1, Hour: 10, Minute: 60, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "minute 60")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__invalidDaysInFebruaryNotLeapYear(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2009, Month: 2, Day: 29, Hour: 10, Minute: 11, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "day 29 for month 2")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__invalidDaysInFebruaryLeapYear(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2016, Month: 2, Day: 30, Hour: 10, Minute: 11, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.ErrorContains(t, err, "day 30 for month 2")
	apiError := err.(*api.ApiError)
	assert.Equal(t, apiError.Code, http.StatusBadRequest)
}

func TestSetSystemDateTime__updateMoreDayInNewMonthThanCurrent(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch {
			case address >= 64045 && address <= 64049:
				from := (address - 64045) * 2
				to := (address + quantity - 64045) * 2
				assert.Check(t, from <= 8, "from %d", from)
				assert.Check(t, to >= 2 && to <= 10, "to %d", from)
				return []byte{0, 10, 0, 11, 0, 14, 0, 2, 7, 229}[from:to], nil // non leap year
			case address == 10198:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 1}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
		WriteSingleRegisterMock: func(address, value uint16) ([]byte, error) {
			return []byte{}, nil
		},
	}
	service := api.NewSystemApiService(mock)
	request := openapi.GetSystemDateTime{Year: 2016, Month: 3, Day: 5, Hour: 9, Minute: 13, AutoDaylightSaving: false}
	_, err := service.SetSystemDateTime(context.TODO(), request)
	assert.NilError(t, err)
	for i, call := range mock.Calls {
		fmt.Printf("% 2d. %v\n", i, call)
	}
	assert.Equal(t, 16, len(mock.Calls))
	assertDeepEqual(t, mock.Calls[3], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(64048), uint16(3)}})    // month
	assertDeepEqual(t, mock.Calls[5], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(64047), uint16(5)}})    // day
	assertDeepEqual(t, mock.Calls[7], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(64049), uint16(2016)}}) // year
	assertDeepEqual(t, mock.Calls[9], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(64045), uint16(9)}})    // hour
	assertDeepEqual(t, mock.Calls[11], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(64046), uint16(13)}})  // minute
}
