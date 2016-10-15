/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This defines the Tag struct, which represents a tag from the DB.
 */

package defs

type Tag struct {
	Name         string `json:"name"`
	CoverImageId uint32 `json:"cover_image_id"`
}
