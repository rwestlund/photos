/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This defines the Photo struct, which represents a photo from the DB.
 */
package defs

import (
	"time"
)

type Photo struct {
	Id           uint32    `json:"id"`
	Filename     string    `json:"filename"`
	Mimetype     string    `json:"mimetype"`
	Size         int64     `json:"size"`
	CreationDate time.Time `json:"creation_date"`
	AuthorId     uint32    `json:"author_id"`
	Caption      string    `json:"caption"`
	/* Fields from other tables. */
	Albums []string `json:"albums"`
}
