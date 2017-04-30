/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file exposes the database interface for albums.
 */

package db

import (
	"database/sql"

	"github.com/rwestlund/photos/defs"
)

// Read an album from SQL rows into a Album object.
func scanAlbum(rows *sql.Rows) (*defs.Album, error) {
	var t defs.Album
	// Because CoverImageID may be null, read into NullInt64 first. The Album
	// object holds a pointer to a uint32 rather than a uint32 directly because
	// this is the only way to make json.Marshal() encode a null when the
	// CoverImageID is not valid.
	var coverimg sql.NullInt64
	var err = rows.Scan(&t.Name, &coverimg, &t.ImageCount)
	if err != nil {
		return nil, err
	}
	if coverimg.Valid {
		var tmp = uint32(coverimg.Int64)
		t.CoverImageID = &tmp
	}
	return &t, nil
}

// FetchAlbums gets a list of all albums in the database.
func FetchAlbums() (*[]defs.Album, error) {
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(`
		SELECT albums.name, albums.cover_image_id,
			COUNT(photo_albums.album_name) AS image_count
		FROM albums
		LEFT JOIN photo_albums ON photo_albums.album_name = albums.name
		GROUP BY albums.name
		ORDER BY albums.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The array we're going to fill. The append() builtin will approximately
	// double the capacity when it needs to reallocate, but we can save some
	// copying by starting at a decent number.
	var albums = make([]defs.Album, 0, 20)
	var album *defs.Album
	// Iterate over rows, reading in each Album as we go.
	for rows.Next() {
		album, err = scanAlbum(rows)
		if err != nil {
			return nil, err
		}
		// Add it to our list.
		albums = append(albums, *album)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &albums, nil
}

// FetchAlbum gets one album from the database.
func FetchAlbum(albumName string) (*defs.Album, error) {
	var rows *sql.Rows
	var err error
	rows, err = DB.Query(`
		SELECT albums.name, albums.cover_image_id,
			COUNT(photo_albums.album_name) AS image_count
		FROM albums
		LEFT JOIN photo_albums ON photo_albums.album_name = albums.name
		WHERE albums.name = $1
		GROUP BY albums.name
		ORDER BY albums.name`, albumName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var album *defs.Album
	album, err = scanAlbum(rows)
	if err != nil {
		return nil, err
	}

	return album, nil
}

// CreateAlbum creates a new album in the database.
func CreateAlbum(album *defs.Album) (*defs.Album, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`INSERT INTO albums (name) VALUES ($1)
                RETURNING name, cover_image_id, 0 AS image_count`, album.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	// Scan it in.
	album, err = scanAlbum(rows)
	if err != nil {
		return nil, err
	}
	return album, nil
}

// UpdateAlbum updates an album in the database.
func UpdateAlbum(name string, album *defs.Album) (*defs.Album, error) {
	var rows *sql.Rows
	var err error
	//TODO some input validation would be nice
	rows, err = DB.Query(`UPDATE albums SET (name, cover_image_id) = ($1, $2)
				WHERE name = $3
                RETURNING name, cover_image_id`,
		album.Name, album.CoverImageID, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Make sure we have a row returned.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	// At this point, we just need to read back the album.
	// TODO replace this with RETURNING
	album, err = FetchAlbum(album.Name)
	return album, err
}
