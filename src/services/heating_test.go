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
	"log"
	"net/http"
	"reflect"
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

func TestSetHeatCurveBySlope__failInvalidSlope(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveBySlopeRequest{
		Slope:       0,
		MinFlowTemp: 10,
		MaxFlowTemp: 150,
	}
	_, err := service.SetHeatCurveBySlope(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "slope")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveBySlope__failInvalidMinFlow(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveBySlopeRequest{
		Slope:       -1,
		MinFlowTemp: 1,
		MaxFlowTemp: 150,
	}
	_, err := service.SetHeatCurveBySlope(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "min flow")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveBySlope__failInvalidMaxFlow(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveBySlopeRequest{
		Slope:       -1,
		MinFlowTemp: 10,
		MaxFlowTemp: 200,
	}
	_, err := service.SetHeatCurveBySlope(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "max flow")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveBySlope__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch address {
			case 11175: // slope
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 17}, nil
			case 11177: // min/max
				assert.Check(t, quantity == 1 || quantity == 2)
				return []byte{0, 33, 0, 66}[0 : quantity*2], nil
			case 11178: // max
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 66}, nil
			case 11400: // temperatures: -30, -15, -5, 0, 5, 15
				assert.Equal(t, quantity, uint16(6))
				return []byte{0, 65, 0, 63, 0, 61, 0, 59, 0, 57, 0, 55}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
		WriteSingleRegisterMock: func(address, value uint16) ([]byte, error) {
			switch address {
			case 11175: // slope
				assert.Equal(t, uint16(18), value)
				return []byte{}, nil
			case 11177: // min
				assert.Equal(t, uint16(30), value)
				return []byte{}, nil
			case 11178: // max
				assert.Equal(t, uint16(70), value)
				return []byte{}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
	}
	service := api.NewHeatingApiService(mock)
	request := openapi.SetHeatCurveBySlopeRequest{
		Slope:       -1.8,
		MinFlowTemp: 30,
		MaxFlowTemp: 70,
	}
	response, err := service.SetHeatCurveBySlope(context.TODO(), 1, request)
	assert.NilError(t, err)
	assert.Check(t, response.Code == http.StatusOK)
	if body, ok := response.Body.(openapi.GetHeatCurveResponse); !ok {
		t.Errorf("Unexpected return type %T\n", body)
	}
	assert.Equal(t, 9, len(mock.Calls))
	for i, call := range mock.Calls {
		fmt.Printf("%d: %v\n", i, call)
	}
	assertDeepEqual(t, mock.Calls[1], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11175), uint16(18)}})
	assertDeepEqual(t, mock.Calls[3], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11177), uint16(30)}})
	assertDeepEqual(t, mock.Calls[5], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11178), uint16(70)}})
}

func TestSetHeatCurveByPoints__failInvalidMinFlow(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveByPointsRequest{
		MinFlowTemp: 1,
		MaxFlowTemp: 150,
		CurvePoints: []openapi.FlowTempPoint{},
	}
	_, err := service.SetHeatCurveByPoints(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "min flow")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveByPoints__failInvalidMaxFlow(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveByPointsRequest{
		MinFlowTemp: 10,
		MaxFlowTemp: 151,
		CurvePoints: []openapi.FlowTempPoint{},
	}
	_, err := service.SetHeatCurveByPoints(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "max flow")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveByPoints__failInvalidCurvePointOutdoorTemp(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveByPointsRequest{
		CurvePoints: []openapi.FlowTempPoint{
			{OutdoorTemp: -7, FlowTemp: 10},
		},
	}
	_, err := service.SetHeatCurveByPoints(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "outdoor temp -7")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveByPoints__failInvalidCurvePointFlowTemp(t *testing.T) {
	mock := &mocks.ClientMock{}
	service := api.NewHeatingApiService(mock)
	values := openapi.SetHeatCurveByPointsRequest{
		CurvePoints: []openapi.FlowTempPoint{
			{OutdoorTemp: 0, FlowTemp: 9},
		},
	}
	_, err := service.SetHeatCurveByPoints(context.TODO(), 1, values)
	assert.ErrorContains(t, err, "9 for flow temp")
	apiErr, ok := err.(*api.ApiError)
	assert.Assert(t, ok, "%T", err)
	assert.Check(t, apiErr.Code == http.StatusBadRequest)
}

func TestSetHeatCurveByPoints__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			switch address {
			case 11175: // slope
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 17}, nil
			case 11177: // min/max
				assert.Check(t, quantity == 1 || quantity == 2)
				return []byte{0, 33, 0, 66}[0 : quantity*2], nil
			case 11178: // max
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 66}, nil
			case 11400: // temperatures: -30, -15, -5, 0, 5, 15
				assert.Check(t, quantity == 1 || quantity == 6)
				return []byte{0, 65, 0, 63, 0, 61, 0, 59, 0, 57, 0, 55}[0 : quantity*2], nil
			case 11401:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 63}, nil
			case 11402:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 61}, nil
			case 11403:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 59}, nil
			case 11404:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 57}, nil
			case 11405:
				assert.Equal(t, uint16(1), quantity)
				return []byte{0, 55}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
		WriteSingleRegisterMock: func(address, value uint16) ([]byte, error) {
			switch address {
			case 11177: // min
				assert.Equal(t, uint16(30), value)
				return []byte{}, nil
			case 11178: // max
				assert.Equal(t, uint16(70), value)
				return []byte{}, nil
			case 11400:
				assert.Equal(t, uint16(10), value)
				return []byte{}, nil
			case 11401:
				assert.Equal(t, uint16(11), value)
				return []byte{}, nil
			case 11402:
				assert.Equal(t, uint16(12), value)
				return []byte{}, nil
			case 11403:
				assert.Equal(t, uint16(13), value)
				return []byte{}, nil
			case 11404:
				assert.Equal(t, uint16(14), value)
				return []byte{}, nil
			case 11405:
				assert.Equal(t, uint16(15), value)
				return []byte{}, nil
			default:
				t.Errorf("Unexpected address %d", address)
				t.FailNow()
				return nil, errors.New("Test failure")
			}
		},
	}
	service := api.NewHeatingApiService(mock)
	request := openapi.SetHeatCurveByPointsRequest{
		MinFlowTemp: 30,
		MaxFlowTemp: 70,
		CurvePoints: []openapi.FlowTempPoint{
			{OutdoorTemp: -30, FlowTemp: 10},
			{OutdoorTemp: -15, FlowTemp: 11},
			{OutdoorTemp: -5, FlowTemp: 12},
			{OutdoorTemp: -0, FlowTemp: 13},
			{OutdoorTemp: 5, FlowTemp: 14},
			{OutdoorTemp: 15, FlowTemp: 15},
		},
	}
	response, err := service.SetHeatCurveByPoints(context.TODO(), 1, request)
	assert.NilError(t, err)
	assert.Check(t, response.Code == http.StatusOK)
	if body, ok := response.Body.(openapi.GetHeatCurveResponse); !ok {
		t.Errorf("Unexpected return type %T\n", body)
	}
	assert.Equal(t, 19, len(mock.Calls))
	for i, call := range mock.Calls {
		fmt.Printf("%d: %v\n", i, call)
	}
	assertDeepEqual(t, mock.Calls[1], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11177), uint16(30)}})
	assertDeepEqual(t, mock.Calls[3], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11178), uint16(70)}})
	assertDeepEqual(t, mock.Calls[5], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11400), uint16(10)}})
	assertDeepEqual(t, mock.Calls[7], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11401), uint16(11)}})
	assertDeepEqual(t, mock.Calls[9], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11402), uint16(12)}})
	assertDeepEqual(t, mock.Calls[11], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11403), uint16(13)}})
	assertDeepEqual(t, mock.Calls[13], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11404), uint16(14)}})
	assertDeepEqual(t, mock.Calls[15], mocks.Call{FuncName: "WriteSingleRegister", Params: []mocks.Param{uint16(11405), uint16(15)}})
}

func assertDeepEqual(t *testing.T, first any, second any) {
	assert.Check(t, reflect.DeepEqual(first, second), "%v != %v\n", first, second)
}
