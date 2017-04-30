/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file exposes the database interface for users.
 */

package db

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/rwestlund/photos/defs"
)

// SQL to select users.
var usersQuery = `
	SELECT users.id, users.email, users.name,
		users.role, users.lastlog, users.creation_date
	FROM users`

// scanUser reads a user from SQL rows into a User object.
func scanUser(rows *sql.Rows) (*defs.User, error) {
	var u defs.User

	// Because lastlog may be null, read into NullTime first. The User object
	// holds a pointer to a time.Time rather than a time.Time directly because
	// this is the only way to make json.Marshal() encode a null when the time
	// is not valid.
	var lastlog pq.NullTime
	// Name may be null, but we're fine converting that to an empty string.
	var name sql.NullString
	var err = rows.Scan(&u.ID, &u.Email, &name, &u.Role, &lastlog,
		&u.CreationDate)
	if err != nil {
		return nil, err
	}
	if lastlog.Valid {
		u.Lastlog = &lastlog.Time
	}
	u.Name = name.String
	return &u, nil
}

// FetchUsers gets all users in the database that match the given filter. The
// query in the filter can match either the name, email, or role.
func FetchUsers(filter *defs.ItemFilter) (*[]defs.User, error) {

	// Hold the dynamically generated portion of our SQL.
	var queryText string
	// Hold all the parameters for our query.
	var params []interface{}

	// Tokenize search string on spaces. Each term must be matched in
	// name or email for a user to be returned.
	var terms = strings.Split(filter.Query, " ")
	// Build and apply having_text.
	for i, term := range terms {
		// Ignore blank terms (comes from leading/trailing spaces).
		if term == "" {
			continue
		}

		if i == 0 {
			queryText += "\n\t WHERE (name ILIKE $"
		} else {
			queryText += " AND (name ILIKE $"
		}
		params = append(params, "%"+term+"%")
		queryText += strconv.Itoa(len(params)) +
			"\n\t\t OR email ILIKE $" +
			strconv.Itoa(len(params)) +
			"\n\t\t OR role ILIKE $" +
			strconv.Itoa(len(params)) + ") "
	}
	queryText += "\n\t ORDER BY lastlog DESC NULLS LAST "

	// Apply count.
	if filter.Count != 0 {
		params = append(params, filter.Count)
		queryText += "\n\t LIMIT $" + strconv.Itoa(len(params))
	}
	// Apply skip.
	if filter.Skip != 0 {
		params = append(params, filter.Count*filter.Skip)
		queryText += "\n\t OFFSET $" + strconv.Itoa(len(params))
	}

	// Run the actual query.
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(usersQuery+queryText, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The array we're going to fill. The append() builtin will approximat
	// double the capacity when it needs to reallocate, but we can save some
	// copying by starting at a decent number.
	var users = make([]defs.User, 0, 20)
	var user *defs.User
	// Iterate over rows, reading in each User as we go.
	for rows.Next() {
		user, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
		// Add it to our list.
		users = append(users, *user)
	}
	return &users, rows.Err()
}

// CreateUser takes a reference to a User and creates it in the database,
// returning the new user. Only User.Email and User.Role are read.
func CreateUser(user *defs.User) (*defs.User, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`INSERT INTO users (email, role) VALUES ($1, $2)
                RETURNING id, email, name, role, lastlog, creation_date`,
		user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	return scanUser(rows)
}

// UpdateUser takes a reference to a User and update it in the database,
// returning fields in the passed object. Only User.ID, User.Email, and
// User.Role are read.
func UpdateUser(id uint32, user *defs.User) (*defs.User, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	// Run one query to update the value.
	rows, err = DB.Query(`UPDATE users SET (email, role) = ($1, $2)
                WHERE id = $3
                RETURNING id, email, name, role, lastlog, creation_date`,
		user.Email, user.Role, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	return scanUser(rows)
}

// DeleteUser deletes a user by id.
func DeleteUser(id uint32) error {
	var err error
	_, err = DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

// UserLogout destroys a login token.
func UserLogout(token string) error {
	var err error
	_, err = DB.Exec(`UPDATE users SET (token, lastlog) =
                (NULL, CURRENT_TIMESTAMP)
            WHERE token = $1`,
		token)
	return err
}

// GoogleLogin records Google login by updating name, token, and lastlog.
func GoogleLogin(email string, name string, token string) (*defs.User, error) {
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(`UPDATE users SET (token, name, lastlog) =
                ($1, $2, CURRENT_TIMESTAMP)
            WHERE email = $3
            RETURNING id, email, name, role, lastlog, creation_date`,
		token, name, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	var user *defs.User
	user, err = scanUser(rows)
	return user, err
}

// FetchUserByToken retrieves a user record based on their authentication token.
func FetchUserByToken(token string) (*defs.User, error) {
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(usersQuery+` WHERE users.token = $1`, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	var user *defs.User
	user, err = scanUser(rows)
	return user, err
}
