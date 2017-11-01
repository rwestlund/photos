/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file builds the actual router from the list of routes.
 */

package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter builds a router by iterating over all routes.
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		// Wrap handler in logger from logger.go.
		var handler http.Handler = Logger(route.handler, route.name)

		router.
			Methods(route.methods...).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)
	}

	// If any client routes fall through to the server, such as during page
	// refresh, send the application back.
	router.
		Methods("GET", "HEAD").
		PathPrefix("/{path:(albums|about|users|uploads)}/").
		Name("path").
		Handler(Logger(ServeIndex, "path"))

	// Add route to handle static files.
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./build/default")))

	return router
}

// ServeIndex manually replies with the index / homepage. This is used to
// support client-side refresh without a hash in the URL.
var ServeIndex = http.HandlerFunc(func(res http.ResponseWriter,
	req *http.Request) {
	log.Println("serving index!")
	http.ServeFile(res, req, "./build/default/index.html")
})
