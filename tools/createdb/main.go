/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 */

package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/rwestlund/photos/config"
)

// main will drop and recreate database and users. This should only be run once
// per deployment, just to initialize things. Run tools/resetdb/main.go next.
func main() {
	var db *sql.DB
	var err error
	// This should be the superuser.
	db, err = sql.Open("postgres",
		"user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: connection params are invalid")
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: failed to connect to the DB")
	}

	log.Println("removing old database")
	wrapSQL(db, "DROP DATABASE IF EXISTS "+config.DatabaseName)
	wrapSQL(db, "DROP USER IF EXISTS "+config.DatabaseUserName)
	log.Println("creating new database")
	wrapSQL(db, "CREATE USER "+config.DatabaseUserName+" WITH LOGIN")
	wrapSQL(db, "CREATE DATABASE "+config.DatabaseName+" WITH OWNER "+
		config.DatabaseUserName)
	log.Println("complete")
}

func wrapSQL(db *sql.DB, s string) {
	_, err := db.Exec(s)
	if err != nil {
		log.Println("error during:", s)
		log.Println(err)
		log.Fatal()
	}
}
