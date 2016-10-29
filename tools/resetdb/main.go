/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * Drop and recreate database objects. Used for testing and creating a new
 * deployment. Must be run after tools/createdb/main.go. Also serves as table
 * documentation.
 */

package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rwestlund/photos/config"
	"log"
)

func main() {
	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", "user="+config.DatabaseUserName+
		" dbname="+config.DatabaseName+" sslmode=disable")
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: connection params are invalid")
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: failed to connect to the DB")
	}

	log.Println("dropping old objects")
	wrap_sql(db, "DROP TABLE IF EXISTS tagged_photos", nil)
	wrap_sql(db, "DROP TABLE IF EXISTS tags", nil)
	wrap_sql(db, "DROP TABLE IF EXISTS photos", nil)
	wrap_sql(db, "DROP TABLE IF EXISTS users", nil)

	log.Println("creating new objects")

	wrap_sql(db, `CREATE TABLE users (
        id              serial PRIMARY KEY,
        email           text NOT NULL,
        name            text,
        role            text NOT NULL,
        token           text,
        creation_date   timestamp WITH TIME ZONE NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
        lastlog         timestamp WITH TIME ZONE
    )`, nil)
	wrap_sql(db, `CREATE TABLE photos (
        id          	serial PRIMARY KEY,
		filename		text NOT NULL,
		mimetype		text NOT NULL,
		size			integer NOT NULL,
        creation_date   timestamp WITH TIME ZONE NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
        author_id   	integer NOT NULL REFERENCES users(id),
		caption			text NOT NULL DEFAULT '',
		image			bytea NOT NULL
    )`, nil)
	wrap_sql(db, `CREATE TABLE tags (
        name            text PRIMARY KEY,
		cover_image_id	integer REFERENCES photos(id) NOT NULL
    )`, nil)
	wrap_sql(db, `CREATE TABLE tagged_photos (
        photo_id	integer REFERENCES photos(id) ON DELETE CASCADE NOT NULL,
        tag_name   	text REFERENCES tags(name) ON DELETE CASCADE NOT NULL,
        UNIQUE (photo_id, tag_name)
    )`, nil)

	log.Println("inserting default values")

	/* Add the default admins listed in the config file to the users table. */
	for _, admin := range config.DefaultAdmins {
		wrap_sql(db, `INSERT INTO users (email, role) VALUES ($1, $2)`,
			[]interface{}{admin, "Admin"})
	}

	log.Println("complete")
}

func wrap_sql(db *sql.DB, s string, params []interface{}) {
	_, err := db.Exec(s, params...)
	if err != nil {
		log.Println("error during:", s)
		log.Println(err)
		log.Fatal()
	}
}
