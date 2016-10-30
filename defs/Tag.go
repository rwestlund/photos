/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This defines the Tag struct, which represents a tag from the DB.
 */

package defs

type Tag struct {
	Name string `json:"name"`
	/* See the comment in db/tags.go:scan_tag() for why this is a pointer. */
	CoverImageId *uint32 `json:"cover_image_id"`
	/* Everything below here is a computed field. */
	/* This is a count of how many images have this tag. */
	ImageCount uint32 `json:"image_count"`
}
