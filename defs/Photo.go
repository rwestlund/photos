/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 */

package defs

import (
	"time"
)

// Photo represents a photo from the DB.
type Photo struct {
	ID           uint32    `json:"id"`
	Filename     string    `json:"filename"`
	Mimetype     string    `json:"mimetype"`
	Size         int64     `json:"size"`
	CreationDate time.Time `json:"creation_date"`
	AuthorID     uint32    `json:"author_id"`
	Caption      string    `json:"caption"`
	/* Fields from other tables. */
	Albums []string `json:"albums"`
}
