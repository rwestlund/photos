/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This is the main file. Run it to launch the application.
 */

package main

import (
	"log"
	"net/http"

	"github.com/rwestlund/photos/config"
	"github.com/rwestlund/photos/db"
	"github.com/rwestlund/photos/router"
)

func main() {
	db.Init()
	/* Create router from routes.go. */
	myRouter := router.NewRouter()
	log.Println("starting server on " + config.ListenAddress)
	log.Fatal(http.ListenAndServe(config.ListenAddress, myRouter))
}
