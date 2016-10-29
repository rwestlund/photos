/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file exposes the database interface for tags.
 */

package db

import (
	"database/sql"
	"github.com/rwestlund/photos/defs"
)

/* Read a tag from SQL rows into a Tag object. */
func scan_tag(rows *sql.Rows) (*defs.Tag, error) {
	var t defs.Tag
	/*
	 * Because CoverImageId may be null, read into NullInt64 first. The Tag
	 * object holds a pointer to a uint32 rather than a uint32 directly because
	 * this is the only way to make json.Marshal() encode a null when the
	 * CoverImageId is not valid.
	 */
	var coverimg sql.NullInt64
	var err error = rows.Scan(&t.Name, &coverimg)
	if err != nil {
		return nil, err
	}
	if coverimg.Valid {
		var tmp uint32 = uint32(coverimg.Int64)
		t.CoverImageId = &tmp
	}
	return &t, nil
}

/* Get a list of all tags in the database. */
func FetchTags() (*[]byte, error) {
	var rows *sql.Rows
	var err error
	/* Return them all in one row. */
	rows, err = DB.Query("SELECT json_agg(name ORDER BY name) FROM tags")
	if err != nil {
		return nil, err
	}
	var tags []byte

	/* In this case, we just want an empty list if nothing was returned. */
	if !rows.Next() {
		return &tags, nil
	}

	/* This is alredy JSON, so just leave it as a []byte. */
	err = rows.Scan(&tags)
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

/* Create a new tag in the database. */
func CreateTag(tag *defs.Tag) (*defs.Tag, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`INSERT INTO tags (name) VALUES ($1)
                RETURNING name, cover_image_id`, tag.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	/* Scan it in. */
	tag, err = scan_tag(rows)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

/* Update a tag in the database. */
func UpdateTag(name string, tag *defs.Tag) (*defs.Tag, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`UPDATE tags SET (name, cover_image_id) = ($1, $2)
				WHERE name = $3
                RETURNING name, cover_image_id`,
		tag.Name, tag.CoverImageId, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	/* Scan it in. */
	tag, err = scan_tag(rows)
	if err != nil {
		return nil, err
	}
	return tag, nil
}
