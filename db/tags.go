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
	var err error = rows.Scan(&t.Name, &coverimg, &t.ImageCount)
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
func FetchTags() (*[]defs.Tag, error) {
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(`
		SELECT tags.name, tags.cover_image_id,
			COUNT(tagged_photos.tag_name) AS image_count
		FROM tags
		LEFT JOIN tagged_photos ON tagged_photos.tag_name = tags.name
		GROUP BY tags.name
		ORDER BY tags.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	/* The array we're going to fill. The append() builtin will approximately
	 * double the capacity when it needs to reallocate, but we can save some
	 * copying by starting at a decent number. */
	var tags = make([]defs.Tag, 0, 20)
	var tag *defs.Tag
	/* Iterate over rows, reading in each Tag as we go. */
	for rows.Next() {
		tag, err = scan_tag(rows)
		if err != nil {
			return nil, err
		}
		/* Add it to our list. */
		tags = append(tags, *tag)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

/* Get one tag from the database. */
func FetchTag(tag_name string) (*defs.Tag, error) {
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(`
		SELECT tags.name, tags.cover_image_id,
			COUNT(tagged_photos.tag_name) AS image_count
		FROM tags
		LEFT JOIN tagged_photos ON tagged_photos.tag_name = tags.name
		WHERE tags.name = $1
		GROUP BY tags.name
		ORDER BY tags.name`, tag_name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var tag *defs.Tag
	tag, err = scan_tag(rows)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

/* Create a new tag in the database. */
func CreateTag(tag *defs.Tag) (*defs.Tag, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`INSERT INTO tags (name) VALUES ($1)
                RETURNING name, cover_image_id, 0 AS image_count`, tag.Name)
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

	/* At this point, we just need to read back the tag. */
	/* TODO replace this with RETURNING */
	tag, err = FetchTag(name)
	return tag, err
}
