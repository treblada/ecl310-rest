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
	"log"
	"net/http"
	"testing"

	"github.com/treblada/ecl310-rest/generated/openapi"
	"github.com/treblada/ecl310-rest/mocks"
	api "github.com/treblada/ecl310-rest/services"
	"gotest.tools/v3/assert"
)

func TestGetHeatCurve__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch address {
			case 11175: // slope
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 17}, nil
			case 11177: // min/max
				assert.Equal(t, uint16(2), quantity)
				return []byte{0, 33, 0, 66}, nil
			case 11400: // temperatures: -30, -15, -5, 0, 5, 15
				assert.Equal(t, uint16(6), quantity)
				return []byte{0, 65, 0, 63, 0, 61, 0, 59, 0, 57, 0, 55}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
	}
	service := api.NewHeatingApiService(mock)
	response, err := service.GetHeatCurve(context.TODO(), 1)
	assert.NilError(t, err)
	assert.Check(t, response.Code == http.StatusOK)
	body := response.Body.(openapi.GetHeatCurveResponse)
	assert.Check(t, body.Slope == -1.7)
	assert.Check(t, body.MinFlowTemp == 33)
	assert.Check(t, body.MaxFlowTemp == 66)
	assert.Check(t, body.CurvePoints[0].OutdoorTemp == -30)
	assert.Check(t, body.CurvePoints[0].FlowTemp == 65)
	assert.Check(t, body.CurvePoints[1].OutdoorTemp == -15)
	assert.Check(t, body.CurvePoints[1].FlowTemp == 63)
	assert.Check(t, body.CurvePoints[2].OutdoorTemp == -5)
	assert.Check(t, body.CurvePoints[2].FlowTemp == 61)
	assert.Check(t, body.CurvePoints[3].OutdoorTemp == -0)
	assert.Check(t, body.CurvePoints[3].FlowTemp == 59)
	assert.Check(t, body.CurvePoints[4].OutdoorTemp == 5)
	assert.Check(t, body.CurvePoints[4].FlowTemp == 57)
	assert.Check(t, body.CurvePoints[5].OutdoorTemp == 15)
	assert.Check(t, body.CurvePoints[5].FlowTemp == 55)
}

func TestGetHeatCurve__failWithCircuit0(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	_, err := service.GetHeatCurve(context.TODO(), 0)
	apiErr, ok := err.(*api.ApiError)
	log.Println(err)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestGetHeatCurve__failWithCircuit4(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	_, err := service.GetHeatCurve(context.TODO(), 4)
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}
