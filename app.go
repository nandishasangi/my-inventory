package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialise(host string, port int, user string, password string, dbname string) error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	app.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	w.WriteHeader(statusCode)
	sendResponse(w, statusCode, error_message)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
	}

	p := product{ID: key}

	err = p.getProduct(app.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusBadRequest, "No Product Found")
		default:
			sendError(w, http.StatusBadRequest, err.Error())
		}
	}
	sendResponse(w, http.StatusOK, p)
}

func (app *App) addProduct(w http.ResponseWriter, r *http.Request) {

	var p product
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Invalid Request Payload")
		return
	}

	err = p.addProduct(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, p)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
	}

	var p product

	err = json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Invalid Request Payload")
		return
	}

	p.ID = key

	err = p.updateProduct(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
	}

	p := product{ID: key}

	err = p.deleteProduct(app.DB)

}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.addProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}
