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
	"log"
	"strconv"
	"strings"

	"github.com/rwestlund/photos/defs"
)

// SQL to select photo
var queryRows = `
	SELECT photos.id, photos.mimetype, photos.size, photos.creation_date,
		photos.author_id, photos.caption, photos.filename, tp.albums
	FROM photos
	LEFT JOIN LATERAL (
		SELECT COALESCE(json_agg(photo_albums.album_name), '[]'::json)
			AS albums
		FROM photo_albums
			WHERE photo_albums.photo_id = photos.id
		) tp ON true`

// Helper function to read Photo out of a sql.Rows object.
func scanPhoto(row *sql.Rows) (*defs.Photo, error) {
	// JSON fields will need special handling.
	var albums string
	// The photo we're going to read in.
	var p defs.Photo

	err := row.Scan(&p.ID, &p.Mimetype, &p.Size, &p.CreationDate, &p.AuthorID,
		&p.Caption, &p.Filename, &albums)
	if err != nil {
		return nil, err
	}
	// Unpack JSON fields.
	e := json.Unmarshal([]byte(albums), &p.Albums)
	if e != nil {
		return nil, e
	}
	return &p, nil
}

// FetchPhotos gets all the photos from the database that match the given
// filter. The query in the filter can match either the caption or the album.
func FetchPhotos(filter *defs.ItemFilter) (*[]defs.Photo, error) {
	_ = log.Println //DEBUG

	// Hold the dynamically generated portion of our SQL.
	var queryText string
	// Hold all the parameters for our query.
	var params []interface{}

	// Tokenize search string on spaces. Each term must be matched in
	// caption or albums for a photo to be returned.
	var terms = strings.Split(filter.Query, " ")
	// Build and apply having_text.
	for i, term := range terms {
		// Ignore blank terms (comes from leading/trailing spaces).
		if term == "" {
			continue
		}

		if i == 0 {
			queryText += "\n\tWHERE (caption ILIKE $"
		} else {
			queryText += " AND (caption ILIKE $"
		}
		params = append(params, "%"+term+"%")
		queryText += strconv.Itoa(len(params)) +
			"\n\t\t OR string_agg(photo_albums.name, ' ') ILIKE $" +
			strconv.Itoa(len(params)) + ") "
	}

	if filter.Album != "" {
		if queryText == "" {
			queryText += "\n\tWHERE"
		} else {
			queryText += "\n\tAND"
		}
		queryText = "\n\t LEFT JOIN photo_albums " +
			"ON photo_albums.photo_id = photos.id " + queryText
		params = append(params, filter.Album)
		queryText += " photo_albums.album_name = $" +
			strconv.Itoa(len(params))
	}

	queryText += "\n\t ORDER BY creation_date DESC"

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
	rows, err = DB.Query(queryRows+queryText, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The array we're going to fill. The append() builtin will approximately
	// double the capacity when it needs to reallocate, but we can save some
	// copying by starting at a decent number.
	var photos = make([]defs.Photo, 0, 20)
	var r *defs.Photo
	// Iterate over rows, reading in each Photo as we go.
	for rows.Next() {
		r, err = scanPhoto(rows)
		if err != nil {
			return nil, err
		}
		// Add it to our list.
		photos = append(photos, *r)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &photos, nil
}

// FetchPhoto gets one photo by id.
func FetchPhoto(id uint32) (*defs.Photo, error) {
	// Read photo from database.
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(queryRows+
		" WHERE photos.id = $1 ", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	var p *defs.Photo
	p, err = scanPhoto(rows)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// FetchPhotoImage gets one photo image by id.
func FetchPhotoImage(id uint32) ([]byte, error) {
	// Read photo from database.
	var rows *sql.Rows
	var err error
	rows, err = DB.Query("SELECT image FROM photos WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
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

// FetchPhotoThumbnail gets one photo thumbnail by id.
func FetchPhotoThumbnail(id uint32) ([]byte, error) {
	// Read photo from database.
	var rows *sql.Rows
	var err error
	rows, err = DB.Query("SELECT thumbnail FROM photos WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
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

// FetchPhotoBigThumbnail gets one photo's big thumbnail by id.
func FetchPhotoBigThumbnail(id uint32) ([]byte, error) {
	// Read photo from database.
	var rows *sql.Rows
	var err error
	rows, err = DB.Query("SELECT big_thumbnail FROM photos WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
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

// CreatePhoto takes a reference to a Photo and creates it in the database,
// returning fields in the passed object.
// TODO this doesn't work yet
func CreatePhoto(photo *defs.Photo, file []byte, thumb []byte,
	bigThumb []byte) (*defs.Photo, error) {
	var rows *sql.Rows
	var err error
	// Start a transaction.
	var tx *sql.Tx
	tx, err = DB.Begin()
	// Implicitly rollback if we exit with an error.
	defer tx.Rollback()

	//TODO some input validation would be nice
	// First we create the photo.
	rows, err = tx.Query(`
			INSERT INTO photos (caption, mimetype, author_id, filename, size,
				image, thumbnail, big_thumbnail)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
		photo.Caption, photo.Mimetype, photo.AuthorID, photo.Filename,
		photo.Size, file, thumb, bigThumb)
	if err != nil {
		return nil, err
	}
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	var id uint32
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}
	// This must be closed before commit; defer doesn't work.
	rows.Close()

	// Insert the new albums.
	// We assume that all the albums exist in the albums table; otherwise we
	// get an error and we'll fail to save.
	var album string
	for _, album = range photo.Albums {
		_, err = tx.Exec(`INSERT INTO photo_albums (photo_id, album_name)
                VALUES ($1, $2)`, id, album)
		if err != nil {
			return nil, err
		}
	}

	// Everything worked, time to commit the transaction.
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// At this point, we just need to read back the new photo.
	photo, err = FetchPhoto(id)
	return photo, err
}

// SavePhoto takes a Photo to save. Only the caption and albums can be modified.
func SavePhoto(photo *defs.Photo) (*defs.Photo, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice

	// Start a transaction.
	var tx *sql.Tx
	tx, err = DB.Begin()
	// Implicitly rollback if we exit with an error.
	defer tx.Rollback()

	// First we update albums. This deleting and then inserting is
	// somewhat wasteful, but it's simple to implement.
	// Use _ and use Exec as in http://go-database-sql.org/modifying.html.
	_, err = tx.Exec("DELETE FROM photo_albums WHERE photo_id = $1", photo.ID)
	if err != nil {
		return nil, err
	}
	// Insert the new albums.
	// We assume that all the albums exist in the albums table; otherwise we
	// get an error and we'll fail to save.
	var album string
	for _, album = range photo.Albums {
		_, err = tx.Exec(`INSERT INTO photo_albums (photo_id, album_name)
                VALUES ($1, $2)`, photo.ID, album)
		if err != nil {
			return nil, err
		}
	}

	// Finally, run the actual query to update the Photo fields.
	rows, err = tx.Query(`UPDATE photos SET caption = $1 WHERE id = $2
		RETURNING id`, photo.Caption, photo.ID)
	if err != nil {
		return nil, err
	}
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	var id uint32
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}
	// This must be closed before commit; defer doesn't work.
	rows.Close()
	// Everything worked, time to commit the transaction.
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	// At this point, we just need to read back the new photo.
	// TODO replace this with RETURNING
	photo, err = FetchPhoto(id)

	return photo, err
}

// DeletePhoto deletes a photo by id.
func DeletePhoto(photoID uint32) error {
	var rows *sql.Rows
	var err error

	rows, err = DB.Query("DELETE FROM photos WHERE id = $1 RETURNING id",
		photoID)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return nil
}
