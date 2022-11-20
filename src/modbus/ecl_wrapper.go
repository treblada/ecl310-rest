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

package wrapper

import "github.com/goburrow/modbus"

type ZeroBasedAddressClientWrapper modbus.Client

/*
The MODbus client wrapper converts the 1-based addresses from the caller to 0-based
addresses used in the modbus library.
*/
type modbusClientWrapper struct {
	ZeroBasedAddressClientWrapper
	client modbus.Client
}

func NewModbusClientWrapper(c modbus.Client) modbusClientWrapper {
	return modbusClientWrapper{
		client: c,
	}
}

func (w *modbusClientWrapper) ReadCoils(address, quantity uint16) (results []byte, err error) {
	return w.client.ReadCoils(address-1, quantity)
}

func (w *modbusClientWrapper) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	return w.client.ReadDiscreteInputs(address-1, quantity)
}

func (w *modbusClientWrapper) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	return w.client.WriteSingleCoil(address-1, value)
}

func (w *modbusClientWrapper) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	return w.client.WriteMultipleCoils(address-1, quantity, value)
}

func (w *modbusClientWrapper) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	return w.client.ReadInputRegisters(address-1, quantity)
}

func (w *modbusClientWrapper) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return w.client.ReadHoldingRegisters(address-1, quantity)
}

func (w *modbusClientWrapper) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	return w.client.WriteSingleRegister(address-1, value)
}

func (w *modbusClientWrapper) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	return w.client.WriteMultipleRegisters(address-1, quantity, value)
}

func (w *modbusClientWrapper) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	return w.client.ReadWriteMultipleRegisters(readAddress-1, readQuantity, writeAddress-1, writeQuantity, value)
}

func (w *modbusClientWrapper) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	return w.client.MaskWriteRegister(address-1, andMask, orMask)
}

func (w *modbusClientWrapper) ReadFIFOQueue(address uint16) (results []byte, err error) {
	return w.client.ReadFIFOQueue(address - 1)
}
