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
	"fmt"
	"net/http"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/treblada/ecl310-rest/generated/openapi"
	"github.com/treblada/ecl310-rest/mocks"
	api "github.com/treblada/ecl310-rest/services"
)

func TestHealth__success(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			return make([]byte, quantity*2), nil
		},
	}
	service := api.NewHealthApiService(mock)
	result, _ := service.GetHealth(context.TODO())
	assert.Equal(t, 200, result.Code)
	bodyContent, ok := result.Body.(openapi.GetHealthResponse)
	assert.Assert(t, ok)
	assert.Equal(t, "OK", bodyContent.Status)
}

func TestHealth__error(t *testing.T) {
	mock := &mocks.ClientMock{
		ReadHoldingRegistersMock: func(address, quantity uint16) ([]byte, error) {
			return nil, fmt.Errorf("Mocked error")
		},
	}
	service := api.NewHealthApiService(mock)
	result, _ := service.GetHealth(context.TODO())
	assert.Equal(t, http.StatusOK, result.Code)
	bodyContent, ok := result.Body.(openapi.GetHealthResponse)
	assert.Assert(t, ok)
	assert.Equal(t, "FAIL", bodyContent.Status)
}
