/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file exposes the database interface for photos.
 */

package db

import (
	"database/sql"
	"encoding/json"
	"github.com/rwestlund/photos/defs"
	"log"
	"strconv"
	"strings"
)

/* SQL to select photos. */
var query_rows string = `
	SELECT photos.id, photos.mimetype, photos.size, photos.creation_date,
		photos.author_id, photos.caption, photos.filename, tp.tags
	FROM photos
	LEFT JOIN LATERAL (
		SELECT COALESCE(json_agg(tagged_photos.tag_name), '[]'::json)
			AS tags
		FROM tagged_photos
			WHERE tagged_photos.photo_id = photos.id
		) tp ON true`

/* Helper function to read Photo out of a sql.Rows object. */
func scan_photo(row *sql.Rows) (*defs.Photo, error) {
	/* JSON fields will need special handling. */
	var tags string
	/* The photo we're going to read in. */
	var p defs.Photo

	err := row.Scan(&p.Id, &p.Mimetype, &p.Size, &p.CreationDate, &p.AuthorId,
		&p.Caption, &p.Filename, &tags)
	if err != nil {
		return nil, err
	}
	/* Unpack JSON fields. */
	e := json.Unmarshal([]byte(tags), &p.Tags)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

/*
 * Fetch all photos from the database that match the given filter. The query
 * in the filter can match either the caption or the tag.
 */
func FetchPhotos(filter *defs.ItemFilter) (*[]defs.Photo, error) {
	_ = log.Println //DEBUG

	/* Hold the dynamically generated portion of our SQL. */
	var query_text string
	/* Hold all the parameters for our query. */
	var params []interface{}

	/* Tokenize search string on spaces. Each term must be matched in the
	 * caption or tags for a photo to be returned.
	 */
	var terms []string = strings.Split(filter.Query, " ")
	/* Build and apply having_text. */
	for i, term := range terms {
		/* Ignore blank terms (comes from leading/trailing spaces). */
		if term == "" {
			continue
		}

		if i == 0 {
			query_text += "\n\tWHERE (caption ILIKE $"
		} else {
			query_text += " AND (caption ILIKE $"
		}
		params = append(params, "%"+term+"%")
		query_text += strconv.Itoa(len(params)) +
			"\n\t\t OR string_agg(tagged_photos.name, ' ') ILIKE $" +
			strconv.Itoa(len(params)) + ") "
	}

	if filter.Tag != "" {
		if query_text == "" {
			query_text += "\n\tWHERE"
		} else {
			query_text += "\n\tAND"
		}
		query_text = "\n\t LEFT JOIN tagged_photos " +
			"ON tagged_photos.photo_id = photos.id " + query_text
		params = append(params, filter.Tag)
		query_text += " tagged_photos.tag_name = $" +
			strconv.Itoa(len(params))
	}

	query_text += "\n\t ORDER BY creation_date DESC"

	/* Apply count. */
	if filter.Count != 0 {
		params = append(params, filter.Count)
		query_text += "\n\t LIMIT $" + strconv.Itoa(len(params))
	}
	/* Apply skip. */
	if filter.Skip != 0 {
		params = append(params, filter.Count*filter.Skip)
		query_text += "\n\t OFFSET $" + strconv.Itoa(len(params))
	}
	/* Run the actual query. */
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(query_rows+query_text, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	/*
	 * The array we're going to fill. The append() builtin will approximately
	 * double the capacity when it needs to reallocate, but we can save some
	 * copying by starting at a decent number.
	 */
	var photos = make([]defs.Photo, 0, 20)
	var r *defs.Photo
	/* Iterate over rows, reading in each Photo as we go. */
	for rows.Next() {
		r, err = scan_photo(rows)
		if err != nil {
			return nil, err
		}
		/* Add it to our list. */
		photos = append(photos, *r)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &photos, nil
}

/* Fetch one photo by id. */
func FetchPhoto(id uint32) (*defs.Photo, error) {
	/* Read photo from database. */
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(query_rows+
		" WHERE photos.id = $1 ", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	/* Scan it in. */
	var p *defs.Photo
	p, err = scan_photo(rows)
	if err != nil {
		return nil, err
	}
	return p, nil
}

/* Fetch one photo image by id. */
func FetchPhotoImage(id uint32) ([]byte, error) {
	/* Read photo from database. */
	var rows *sql.Rows
	var err error
	rows, err = DB.Query("SELECT image FROM photos WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var image []byte
	err = rows.Scan(&image)
	if err != nil {
		return nil, err
	}
	return image, nil
}

/*
 * Take a reference to a Photo and create it in the database, returning fields
 * in the passed object.
 * TODO this doesn't work yet
 */
func CreatePhoto(photo *defs.Photo, file []byte) (*defs.Photo, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`
			INSERT INTO photos (caption, mimetype, author_id, filename, size, image)
            VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id`,
		photo.Caption, photo.Mimetype, photo.AuthorId, photo.Filename,
		photo.Size, file)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	/* Scan it in. */
	var id uint32
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}
	/* At this point, we just need to read back the new photo. */
	photo, err = FetchPhoto(id)

	return photo, err
}

/*
 * Take a Photo to save. Only the caption and tags can be modified.
 */
func SavePhoto(photo *defs.Photo) (*defs.Photo, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice

	/* Start a transaction. */
	var tx *sql.Tx
	tx, err = DB.Begin()
	/* Implicitly rollback if we exit with an error. */
	defer tx.Rollback()

	/*
	 * First we update tags. This deleting and then inserting is
	 * somewhat wasteful, but it's simple to implement.
	 */
	/* Use _ and use Exec as in http://go-database-sql.org/modifying.html. */
	_, err = tx.Exec("DELETE FROM tagged_photos WHERE photo_id = $1", photo.Id)
	if err != nil {
		return nil, err
	}
	/* Insert the new tags. */
	/* We assume that all the tags exist in the tags table; otherwise we'll
	 * get an error and we'll fail to save. */
	var tag string
	for _, tag = range photo.Tags {
		_, err = tx.Exec(`INSERT INTO tagged_photos (photo_id, tag_name)
                VALUES ($1, $2)`, photo.Id, tag)
		if err != nil {
			return nil, err
		}
	}

	/* Finally, run the actual query to update the Photo fields. */
	rows, err = tx.Query(`UPDATE photos SET caption = $1 WHERE id = $2
		RETURNING id`, photo.Caption, photo.Id)
	if err != nil {
		return nil, err
	}
	/* Make sure we have a row returned. */
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	/* Scan it in. */
	var id uint32
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}
	/* This must be closed before commit; defer doesn't work. */
	rows.Close()
	/* Everything worked, time to commit the transaction. */
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	/* At this point, we just need to read back the new photo. */
	/* TODO replace this with RETURNING */
	photo, err = FetchPhoto(id)

	return photo, err
}

/*
 * Delete a Photo by id.
 */
func DeletePhoto(photo_id uint32) error {
	var rows *sql.Rows
	var err error

	rows, err = DB.Query("DELETE FROM photos WHERE id = $1 RETURNING id",
		photo_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return nil
}
