/*
 * Copyright (c) 2016-2017, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This is the main file. Run it to launch the application.
 */

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rwestlund/photos/db"
	"github.com/rwestlund/photos/defs"
	"github.com/rwestlund/photos/router"
	"github.com/rwestlund/photos/util"
)

func main() {
	flag.Parse()
	// Determine which command to execute.
	switch flag.Arg(0) {
	case "launch":
		launch()
	case "createdb":
		createdb()
	default:
		usage()
	}
}

func usage() {
	os.Stderr.WriteString(`Photos takes one of the following subcommands:
	launch 		Launch the server.
	createdb 	Create the database tables.
	`)
}

func launch() {
	var config, err = defs.ReadConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	db.Init(config)
	myRouter := router.NewRouter(config)
	log.Println("starting server on " + config.ListenAddress)
	log.Fatal(http.ListenAndServe(config.ListenAddress, myRouter))
}

func createdb() {
	var config, err = defs.ReadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	// Give the user a chance to cancel.
	os.Stderr.WriteString("\n\tWARNING: This will destroy all data in the " +
		config.DatabaseName +
		" database.\n\t\tIf this is a mistake, hit ^C NOW!\n")
	time.Sleep(5 * time.Second)

	err = util.InitDB(config)
	if err != nil {
		log.Fatal(err)
	}
	err = util.CreateDB(config)
	if err != nil {
		log.Fatal(err)
	}

}
