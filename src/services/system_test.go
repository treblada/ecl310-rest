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
	"testing"

	"github.com/treblada/ecl310-rest/generated/openapi"
	"github.com/treblada/ecl310-rest/mocks"
	api "github.com/treblada/ecl310-rest/services"
	"gotest.tools/v3/assert"
)

func TestGetSystemInfo_success(t *testing.T) {
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
	body, ok := response.Body.(openapi.GetSystemInfoResponse)
	assert.Assert(t, ok)
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
	response, err := service.GetSystemInfo(context.TODO())
	assert.NilError(t, err)
	assert.Equal(t, 502, response.Code)
	_, ok := response.Body.(string)
	assert.Assert(t, ok)
}
