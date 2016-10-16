/*
 * Copyright (c) 2016, Randy Westlund. All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file contains HTTP handlers for the application.
 */

package router

import (
    "strconv"
    "log"
    "net/url"
    "net/http"
    "encoding/json"
    "database/sql"
    "github.com/rwestlund/photos/defs"
    "github.com/rwestlund/photos/db"
    "github.com/gorilla/mux"
)

/*
 * Take a url.URL object (from req.URL) and fill an ItemFilter.
 */
func build_item_filter(url *url.URL) *defs.ItemFilter {
    /* We can ignore the error because count=0 means disabled. */
    var bigcount uint64
    bigcount, _ = strconv.ParseUint(url.Query().Get("count"), 10, 32)
    var bigskip uint64
    bigskip, _ = strconv.ParseUint(url.Query().Get("skip"), 10, 32)
    /* Build ItemFilter from query params. */
    var filter defs.ItemFilter = defs.ItemFilter{
        Query: url.Query().Get("query"),
        Count: uint32(bigcount),
        Skip: uint32(bigskip),
    }
    return &filter
}

/*
 * Request a list of photos.
 * GET /api/photos
 */
func handle_photos(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "application/json; charset=UTF-8")

    var filter = build_item_filter(req.URL);

    var photos *[]defs.Photo
    var err error
    photos, err = db.FetchPhotos(filter)
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
    /* If we made it here, send good response. */
    res.Write(j)
}

/*
 * Create a new photo or update an existing one.
 * POST /api/photos, PUT /api/photos/4
 */
func handle_put_or_post_photo(res http.ResponseWriter, req *http.Request) {
    /* Access control. */
    var usr *defs.User
    var err error
    usr, err = check_auth(res, req)
    if err != nil {
        res.WriteHeader(500)
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

    /* Decode body. */
    var photo defs.Photo
    err = json.NewDecoder(req.Body).Decode(&photo)
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }

    var new_photo *defs.Photo

    /* Update it. */
    if req.Method == "PUT" {
        new_photo, err = db.SavePhoto(&photo)
        if err == sql.ErrNoRows {
            res.WriteHeader(404)
            return
        }
    /* Create it. */
    } else {
        /* Fill in the currently logged-in user as the author. */
        photo.AuthorId = usr.Id
        new_photo, err = db.CreatePhoto(&photo)
    }

    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }

    /* Send it back. */
    j, e := json.Marshal(new_photo)
    if e != nil {
        log.Println(e)
        res.WriteHeader(500)
        return
    }
    /* If we made it here, send good response. */
    res.Write(j)
}

/*
 * Request a specific photo.
 * GET /api/photo/3
 */
func handle_photo(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "application/json; charset=UTF-8")

    /* Get id parameter. */
    var params map[string]string = mux.Vars(req)
    bigid, err := strconv.ParseUint(params["id"], 10, 32)
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }
    var id uint32 = uint32(bigid)

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

    /* If we made it here, send good response. */
    res.Write(j)
}

/*
 * Delete a photo by id.
 * DELETE /api/photos/4
 */
func handle_delete_photo(res http.ResponseWriter, req *http.Request) {
    /* Access control. */
    var usr *defs.User
    var err error
    usr, err = check_auth(res, req)
    if err != nil {
        res.WriteHeader(500)
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

    /* Get id parameter. */
    var params map[string]string = mux.Vars(req)
    var bigid uint64
    bigid, err = strconv.ParseUint(params["id"], 10, 32)
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }
    var photo_id uint32 = uint32(bigid)

    err = db.DeletePhoto(photo_id)
    if err == sql.ErrNoRows {
        res.WriteHeader(404)
        return
    }
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }
    /* If we made it here, send good response. */
    res.WriteHeader(200)
}

/*
 * Request a list of users.
 * GET /api/users
 */
func handle_users(res http.ResponseWriter, req *http.Request) {
    /* Access control. */
    var usr *defs.User
    var err error
    usr, err = check_auth(res, req)
    if err != nil {
        res.WriteHeader(500)
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

    var filter = build_item_filter(req.URL);

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
    /* If we made it here, send good response. */
    res.Write(j)
}


/*
 * Receive a new user to create.
 * POST /api/users or PUT /api/users/4
 * Example: { email: ..., role: ... }
 */
func handle_post_or_put_user(res http.ResponseWriter, req *http.Request) {
    /* Access control. */
    var usr *defs.User
    var err error
    usr, err = check_auth(res, req)
    if err != nil {
        res.WriteHeader(500)
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

    /* Decode body. */
    var user defs.User
    err = json.NewDecoder(req.Body).Decode(&user)
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }

    var new_user *defs.User
    /* Update a user in the database. */
    if req.Method == "PUT" {
        /* Get id parameter. */
        var params map[string]string = mux.Vars(req)
        bigid, err := strconv.ParseUint(params["id"], 10, 32)
        if err != nil {
            log.Println(err)
            res.WriteHeader(400)
            return
        }
        var id uint32 = uint32(bigid)

        new_user, err = db.UpdateUser(id, &user)
    /* Create new user in DB. */
    } else {
        new_user, err = db.CreateUser(&user)
    }

    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }

    /* Send it back. */
    j, e := json.Marshal(new_user)
    if e != nil {
        log.Println(e)
        res.WriteHeader(500)
        return
    }
    /* If we made it here, send good response. */
    res.Write(j)
}


/*
 * Delete a user by id.
 * DELETE /api/users/4
 */
func handle_delete_user(res http.ResponseWriter, req *http.Request) {
    /* Access control. */
    var usr *defs.User
    var err error
    usr, err = check_auth(res, req)
    if err != nil {
        res.WriteHeader(500)
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

    /* Get id parameter. */
    var params map[string]string = mux.Vars(req)
    var bigid uint64
    bigid, err = strconv.ParseUint(params["id"], 10, 32)
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }
    var id uint32 = uint32(bigid)

    err = db.DeleteUser(id)

    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        return
    }

    /* If we made it here, send good response. */
    res.WriteHeader(200)
}

func handle_get_tags(res http.ResponseWriter, req *http.Request) {
    var tags *[]byte
    var err error
    tags, err = db.FetchTags()
    if err != nil {
        log.Println(err)
        res.WriteHeader(500)
        return
    }
    res.Write(*tags)
}


