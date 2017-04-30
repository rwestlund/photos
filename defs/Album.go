/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 */

package defs

// Album represents an album in the DB.
type Album struct {
	Name string `json:"name"`
	// See the comment in db/albums.go:scan_album() for why this is a pointer.
	CoverImageID *uint32 `json:"cover_image_id"`
	// Everything below here is a computed field.
	// This is a count of how many images have this album.
	ImageCount uint32 `json:"image_count"`
}
