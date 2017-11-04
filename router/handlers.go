/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file contains HTTP handlers for the application.
 */

package router

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/rwestlund/photos/db"
	"github.com/rwestlund/photos/defs"
)

// buildItemFilder takes a url.URL object from req.URL and fills an ItemFilter.
func buildItemFilter(url *url.URL) *defs.ItemFilter {
	// We can ignore the error because count=0 means disabled.
	var bigcount, _ = strconv.ParseUint(url.Query().Get("count"), 10, 32)
	var bigskip, _ = strconv.ParseUint(url.Query().Get("skip"), 10, 32)
	// Build ItemFilter from query params.
	var filter = defs.ItemFilter{
		Query: url.Query().Get("query"),
		Count: uint32(bigcount),
		Skip:  uint32(bigskip),
		Album: url.Query().Get("album"),
	}
	return &filter
}

// handlePhotos requests a list of photos.
// GET /api/photos
func handlePhotos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var filter = buildItemFilter(req.URL)

	var photos, err = db.FetchPhotos(filter)
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}
	j, e := json.Marshal(photos)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

// handlePutPhoto updates an existing photo.
// PUT /api/photos/4
func handlePutPhoto(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Decode body.
	var photo defs.Photo
	err = json.NewDecoder(req.Body).Decode(&photo)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var newPhoto *defs.Photo

	// Update it.
	newPhoto, err = db.SavePhoto(&photo)
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	}
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Send it back.
	j, e := json.Marshal(newPhoto)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

// handlePostPhoto creates a new photo.
// POST /api/photos
func handlePostPhoto(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Hold the first 200MB in RAM; the rest goes to temporary files.
	err = req.ParseMultipartForm(200 * 1024 * 1024)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var file multipart.File
	var header *multipart.FileHeader
	file, header, err = req.FormFile("file")
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Decode body.
	var photo = defs.Photo{
		Filename: header.Filename,
		Mimetype: header.Header.Get("Content-Type"),
	}

	// Full size image.
	var photoBuff bytes.Buffer
	photo.Size, err = photoBuff.ReadFrom(file)

	// Create thumbnails.
	var img image.Image
	var imgType string
	img, imgType, err = image.Decode(&photoBuff)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	// Create small thumbnail.
	var thumb = resize.Thumbnail(800, 800, img, resize.Lanczos3)

	// Create big thumbnail.
	var bigThumb = resize.Thumbnail(1600, 1600, img, resize.Lanczos3)

	// Put thumbnail into a byte arrays for the database.
	var thumbBuff bytes.Buffer
	if imgType == "jpeg" {
		err = jpeg.Encode(&thumbBuff, thumb, nil)
	} else if imgType == "png" {
		err = png.Encode(&thumbBuff, thumb)
	} else {
		log.Println("Unsupported image type: " + imgType)
		res.WriteHeader(400)
		return
	}
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	// Put big thumbnail into a byte array for the database.
	var bigThumbBuff bytes.Buffer
	if imgType == "jpeg" {
		err = jpeg.Encode(&bigThumbBuff, bigThumb, nil)
	} else if imgType == "png" {
		err = png.Encode(&bigThumbBuff, bigThumb)
	} else {
		log.Println("Unsupported image type: " + imgType)
		res.WriteHeader(400)
		return
	}
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Add albums
	var albumsString = req.FormValue("albums")
	err = json.Unmarshal([]byte(albumsString), &photo.Albums)
	if err != nil {
		log.Println(err)
		log.Println(photo.Albums)
	}

	var newPhoto *defs.Photo
	// Fill in the currently logged-in user as the author.
	photo.AuthorID = usr.ID
	newPhoto, err = db.CreatePhoto(&photo, photoBuff.Bytes(),
		thumbBuff.Bytes(), bigThumbBuff.Bytes())

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Send it back.
	j, e := json.Marshal(newPhoto)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

// handlePhoto requests a specific photo.
// GET /api/photo/3
func handlePhoto(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get id parameter.
	bigid, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	var id = uint32(bigid)

	var photo *defs.Photo
	photo, err = db.FetchPhoto(id)
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	} else if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	j, e := json.Marshal(photo)
	if e != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	// If we made it here, send good response.
	res.Write(j)
}

// handlePhotoImage requests a specific photo image.
// GET /api/photo/3/image
func handlePhotoImage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/binary")

	// Get id parameter.
	bigid, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	var id = uint32(bigid)

	var image []byte
	image, err = db.FetchPhotoImage(id)
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	} else if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}

	// If we made it here, send good response.
	res.Write(image)
}

// handlePhotoThumbnail requests a specific photo.
// GET /api/photo/3/thumbnail
func handlePhotoThumbnail(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/binary")

	// Get id parameter.
	bigid, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	var id = uint32(bigid)

	var image []byte
	image, err = db.FetchPhotoThumbnail(id)
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	} else if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}

	// If we made it here, send good response.
	res.Write(image)
}

// handlePhotoBigThumbnail requests a specific photo's big thumbnail.
// GET /api/photo/3/big_thumbnail
func handlePhotoBigThumbnail(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/binary")

	// Get id parameter.
	bigid, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	var id = uint32(bigid)

	var image []byte
	image, err = db.FetchPhotoBigThumbnail(id)
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	} else if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}

	// If we made it here, send good response.
	res.Write(image)
}

// handleDeletePhoto deletes a photo by id.
// DELETE /api/photos/4
func handleDeletePhoto(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get id parameter.
	var bigid uint64
	bigid, err = strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	var photoID = uint32(bigid)

	err = db.DeletePhoto(photoID)
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	}
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	// If we made it here, send good response.
	res.WriteHeader(200)
}

// handleUsers requests a list of users.
// GET /api/users
func handleUsers(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var filter = buildItemFilter(req.URL)

	var users *[]defs.User
	users, err = db.FetchUsers(filter)
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}
	j, e := json.Marshal(users)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

// handlePostOrPutUser receive a new user to create.
// POST /api/users or PUT /api/users/4
// Example: { email: ..., role: ... }
func handlePostOrPutUser(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Decode body.
	var user defs.User
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var newUser *defs.User
	// Update a user in the database.
	if req.Method == "PUT" {
		// Get id parameter.
		bigid, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
		if err != nil {
			log.Println(err)
			res.WriteHeader(400)
			return
		}
		var id = uint32(bigid)

		newUser, err = db.UpdateUser(id, &user)
		// Create new user in DB.
	} else {
		newUser, err = db.CreateUser(&user)
	}

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Send it back.
	j, e := json.Marshal(newUser)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

// handleDeleteUser deletes a user by id.
// DELETE /api/users/4
func handleDeleteUser(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get id parameter.
	var bigid uint64
	bigid, err = strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	var id = uint32(bigid)

	err = db.DeleteUser(id)

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// If we made it here, send good response.
	res.WriteHeader(200)
}

func handleGetAlbums(res http.ResponseWriter, req *http.Request) {
	var albums, err = db.FetchAlbums()
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}
	j, e := json.Marshal(albums)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

func handleGetAlbum(res http.ResponseWriter, req *http.Request) {
	// Get name parameter.
	var album, err = db.FetchAlbum(mux.Vars(req)["albumName"])
	if err == sql.ErrNoRows {
		res.WriteHeader(404)
		return
	}
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}
	j, e := json.Marshal(album)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

func handlePostAlbums(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Decode body.
	var album defs.Album
	err = json.NewDecoder(req.Body).Decode(&album)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var newAlbum *defs.Album
	// Create new album in DB.
	newAlbum, err = db.CreateAlbum(&album)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Send it back.
	j, e := json.Marshal(newAlbum)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)
}

// handlePutAlbums
func handlePutAlbums(res http.ResponseWriter, req *http.Request) {
	// Access control.
	var usr, err = checkAuth(res, req)
	if err != nil {
		res.WriteHeader(500)
		log.Println(err)
		return
	}
	if usr == nil {
		res.WriteHeader(401)
		return
	}
	if usr.Role != "Admin" {
		res.WriteHeader(403)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Decode body.
	var album defs.Album
	err = json.NewDecoder(req.Body).Decode(&album)
	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var newAlbum *defs.Album
	// Update an album in the database.
	// Get name parameter.
	var params = mux.Vars(req)
	newAlbum, err = db.UpdateAlbum(params["name"], &album)

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	// Send it back.
	j, e := json.Marshal(newAlbum)
	if e != nil {
		log.Println(e)
		res.WriteHeader(500)
		return
	}
	// If we made it here, send good response.
	res.Write(j)

}
