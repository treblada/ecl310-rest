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

/*
Unfortunatelly I have to sneak this file among the generated files to be able to pass
a customer error handler to controllers.
*/
package openapi

func NewHealthApiControllerWithErrorHandler(s HealthApiServicer, h ErrorHandler, opts ...HealthApiOption) Router {
	controller := &HealthApiController{
		service:      s,
		errorHandler: h,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

func NewSystemApiControllerWithErrorHandler(s SystemApiServicer, h ErrorHandler, opts ...SystemApiOption) Router {
	controller := &SystemApiController{
		service:      s,
		errorHandler: h,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

func NewHeatingApiControllerWithErrorHandler(s HeatingApiServicer, h ErrorHandler, opts ...HeatingApiOption) Router {
	controller := &HeatingApiController{
		service:      s,
		errorHandler: h,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}
