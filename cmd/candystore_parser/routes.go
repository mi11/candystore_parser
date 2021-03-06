package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/customers", app.getAllCustomers)
	router.HandlerFunc(http.MethodGet, "/top", app.getTopCustomers)

	return router
}
