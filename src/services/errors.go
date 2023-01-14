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

import "fmt"

type ApiError struct {
	error
	Code    int
	Message string
	Cause   error
}

func (err *ApiError) Error() string {
	var causeMessage string = ""
	if err.Cause != nil {
		causeMessage = fmt.Sprintf("; %s", err.Cause.Error())
	}
	return fmt.Sprintf("HTTP %d; %s%s", err.Code, err.Message, causeMessage)
}

func NewApiError(code int, message string, cause error) *ApiError {
	return &ApiError{Code: code, Message: message, Cause: cause}
}
