/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 */

package defs

import (
	"time"
)

// User represents a user from the DB.
type User struct {
	ID    uint32 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	// See the comment in db/users.go:scan_user() for why this is a pointer.
	Lastlog      *time.Time `json:"lastlog"`
	CreationDate time.Time  `json:"creation_date"`
}
