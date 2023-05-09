package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requireAuthenticatedUser(app.healthcheckHandler))

	router.HandlerFunc(http.MethodPost, "/v1/insurance", app.requireLocalConnection(app.createInsuranceHandler))
	router.HandlerFunc(http.MethodGet, "/v1/insurance/:id", app.showInsuranceHandler)
	router.HandlerFunc(http.MethodGet, "/v1/insurance", app.listInsuranceHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/insurance/:id", app.requireLocalConnection(app.updateInsuranceHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/insurance/:id", app.requireLocalConnection(app.deleteInsuranceHandler))

	router.HandlerFunc(http.MethodPost, "/v1/order", app.createOrderHandler)
	router.HandlerFunc(http.MethodGet, "/v1/order/:id", app.showOrderHandler)
	router.HandlerFunc(http.MethodGet, "/v1/order", app.listOrdersHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/order/:id", app.updateOrderHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/order/:id", app.deleteOrderHandler)

	router.HandlerFunc(http.MethodPost, "/v1/order_properties", app.requireLocalConnection(app.createPropertyHandler))
	router.HandlerFunc(http.MethodGet, "/v1/order_properties/:id", app.showPropertyHandler)
	router.HandlerFunc(http.MethodGet, "/v1/order_properties", app.listPropertiesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/order_properties/:id", app.requireLocalConnection(app.updatePropertyHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/order_properties/:id", app.requireLocalConnection(app.deletePropertyHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/admin/healthcheck", app.requireLocalConnection(app.healthcheckHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/query", app.requireLocalConnection(app.customQueryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/tables", app.requireLocalConnection(app.listTablesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/table_description", app.requireLocalConnection(app.describeTableHandler))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
