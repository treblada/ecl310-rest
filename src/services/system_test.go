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
			assert.Check(t, quantity == 3)
			switch address {
			case 4201:
				return []byte{0, 2, 0, 3, 0, 4}, nil
			case 4211:
				return []byte{0, 1, 0, 2, 0, 3}, nil
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
