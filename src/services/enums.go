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

type CircuitMode uint16

const (
	Manual CircuitMode = iota
	Scheduled
	ConstantComfortTemp
	ConstantSetbackTemp
	FrostProtection
)

var circuitModeNames = []string{"MANUAL", "SCHEDULED", "CONSTANT_COMFORT_TEMP", "CONSTANT_SETBACK_TEMP", "FROST_PROTECTION"}

// var circuitModeEnums = []CircuitMode{Manual, Scheduled, ConstantComfortTemp, ConstantSetbackTemp, FrostProtection}

func (m CircuitMode) String() string {
	return circuitModeNames[m]
}

func GetCircuitMode(i uint16) CircuitMode {
	return CircuitMode(i)
}

type CircuitState uint16

const (
	Setback CircuitState = iota
	PreComfort
	Comfort
	PreSetback
)

var circuitStateNames = []string{"SETBACK", "PRE_COMFORT", "COMFORT", "PRE_SETBACK"}

func (s CircuitState) String() string {
	return circuitStateNames[s]
}

func GetCircuitState(i uint16) CircuitState {
	return CircuitState(i)
}

type AddressType uint16

const (
	Static AddressType = iota
	DHCP
)

var addressTypeNames = []string{"STATIC", "DHCP"}

func (t AddressType) String() string {
	return addressTypeNames[t]
}

func GetAddressType(i uint16) AddressType {
	return AddressType(i)
}
