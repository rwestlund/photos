// Copyright (c) 2017, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package util

import (
	"database/sql"
	"log"

	"github.com/rwestlund/photos/defs"
)

// InitDB drops and recreates the bare datebase.
func InitDB(config *defs.Config) error {
	// Connect as the superuser.
	var db, err = sql.Open("postgres",
		"user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	var statements = [][]string{
		{"Dropping old database and roles...",
			"DROP DATABASE IF EXISTS " + config.DatabaseName},
		{"", "DROP ROLE IF EXISTS " + config.DatabaseUserName},
		{"Creating new database and roles...",
			"CREATE ROLE " + config.DatabaseUserName + " WITH LOGIN"},
		{"", "CREATE DATABASE " + config.DatabaseName +
			" WITH OWNER " + config.DatabaseUserName},
	}
	for _, s := range statements {
		if s[0] != "" {
			log.Println(s[0])
		}
		_, err = db.Exec(s[1])
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDB creates the database tables.
func CreateDB(config *defs.Config) error {
	var db, err = sql.Open("postgres", "user="+config.DatabaseUserName+
		" dbname="+config.DatabaseName+" sslmode=disable")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	var statements = [][]string{
		{"Creating table users...", `CREATE TABLE users (
			id              serial PRIMARY KEY,
			email           text NOT NULL,
			name            text,
			role            text NOT NULL,
			token           text,
			creation_date   timestamp WITH TIME ZONE NOT NULL
								DEFAULT CURRENT_TIMESTAMP,
			lastlog         timestamp WITH TIME ZONE)`},
		{"Creating table photos...", `CREATE TABLE photos (
			id          	serial PRIMARY KEY,
			filename		text NOT NULL,
			mimetype		text NOT NULL,
			size			integer NOT NULL,
			creation_date   timestamp WITH TIME ZONE NOT NULL
								DEFAULT CURRENT_TIMESTAMP,
			author_id   	integer NOT NULL REFERENCES users(id),
			caption			text NOT NULL DEFAULT '',
			image			bytea NOT NULL,
			thumbnail		bytea NOT NULL,
			big_thumbnail	bytea NOT NULL)`},
		{"Creating table albums...", `CREATE TABLE albums (
			name            text PRIMARY KEY,
			cover_image_id	integer REFERENCES photos(id))`},
		{"Creating table photo_albums...", `CREATE TABLE photo_albums (
			photo_id	integer REFERENCES photos(id) ON DELETE CASCADE NOT NULL,
			album_name   	text REFERENCES albums(name) ON DELETE CASCADE
							ON UPDATE CASCADE NOT NULL,
			UNIQUE (photo_id, album_name))`},
	}
	for _, s := range statements {
		if s[0] != "" {
			log.Println(s[0])
		}
		_, err = db.Exec(s[1])
		if err != nil {
			return err
		}
	}

	// Add the default admins listed in the config file to the users table.
	log.Println("Inserting default admins...")
	for _, admin := range config.DefaultAdmins {
		_, err = db.Exec(`INSERT INTO users (email, role) VALUES ($1, $2)`,
			admin, "Admin")
		if err != nil {
			return err
		}
	}

	log.Println("Database setup complete!")
	return nil
}
