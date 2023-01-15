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

package mocks

import (
	wrapper "github.com/treblada/ecl310-rest/modbus"
)

type Param interface{}

type Call struct {
	FuncName string
	Params   []Param
}

type ClientMock struct {
	wrapper.ZeroBasedAddressClientWrapper
	ReadHoldingRegistersMock func(address, quantity uint16) (results []byte, err error)
	WriteSingleRegisterMock  func(address, value uint16) (results []byte, err error)
	Calls                    []Call
}

func (c *ClientMock) registerCall(funcName string, params ...Param) {
	c.Calls = append(c.Calls, Call{FuncName: funcName, Params: params})
}

func (c *ClientMock) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	c.registerCall("ReadHoldingRegisters", address, quantity)
	return c.ReadHoldingRegistersMock(address, quantity)
}

func (c *ClientMock) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	c.registerCall("WriteSingleRegister", address, value)
	return c.WriteSingleRegisterMock(address, value)
}
