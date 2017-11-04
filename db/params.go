/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file connects to the database and exposes the handle to the other DB
 * files.
 */

package db

import (
	"database/sql"
	"log"

	// Import the postgres driver.
	_ "github.com/lib/pq"
	"github.com/rwestlund/photos/defs"
)

// DB is the db handle for this package.
var DB *sql.DB

// Init connects to the database.
func Init(config *defs.Config) {
	// Connect to database.
	var err error
	DB, err = sql.Open("postgres", "user="+config.DatabaseUserName+
		" dbname="+config.DatabaseName+" sslmode=disable")
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: connection params are invalid")
	}
	err = DB.Ping()
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: failed to connect to the DB")
	}
}
