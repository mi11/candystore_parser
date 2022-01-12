package main

import "net/http"

func (app *application) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := app.models.ExtractCustomers()
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, customers, "customers")
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
	}
}

func (app *application) getTopCustomers(w http.ResponseWriter, r *http.Request) {
	topCustomers, err := app.models.GetTopCustomers()
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, topCustomers, "topCustomers")
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
	}
}
