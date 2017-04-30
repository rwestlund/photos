/*
 * Copyright (c) 2016, Randy Westlund and Jacqueline Kory Westlund.
 * All rights reserved.
 * This code is under the BSD-2-Clause license.
 *
 * This file defines the application's routes, mapping them to handlers.
 */

package router

import (
	"net/http"
)

// Route contains the information needed for each Route.
type Route struct {
	name    string
	methods []string
	pattern string
	handler http.HandlerFunc
}

// Routes is a list of routes.
type Routes []Route

// Define the actual routes here.
var routes = Routes{
	Route{
		"auth",
		[]string{"GET"},
		"/auth/google/login",
		oauthRedirect,
	},
	Route{
		"auth",
		[]string{"GET"},
		"/auth/google/return",
		handleOauthCallback,
	},
	Route{
		"logout",
		[]string{"GET"},
		"/auth/logout",
		handleLogout,
	},
	Route{
		"photo",
		[]string{"GET", "HEAD"},
		"/api/photos/{id:[0-9]+}",
		handlePhoto,
	},
	Route{
		"photo",
		[]string{"GET", "HEAD"},
		"/api/photos/{id:[0-9]+}/image",
		handlePhotoImage,
	},
	Route{
		"photo",
		[]string{"GET", "HEAD"},
		"/api/photos/{id:[0-9]+}/thumbnail",
		handlePhotoThumbnail,
	},
	Route{
		"photo",
		[]string{"GET", "HEAD"},
		"/api/photos/{id:[0-9]+}/big_thumbnail",
		handlePhotoBigThumbnail,
	},
	Route{
		"photos",
		[]string{"GET", "HEAD"},
		"/api/photos",
		handlePhotos,
	},
	Route{
		"users",
		[]string{"GET", "HEAD"},
		"/api/users",
		handleUsers,
	},
	Route{
		"photos",
		[]string{"POST"},
		"/api/photos",
		handlePostPhoto,
	},
	Route{
		"photos",
		[]string{"PUT"},
		"/api/photos/{id:[0-9]+}",
		handlePutPhoto,
	},
	Route{
		"photos",
		[]string{"DELETE"},
		"/api/photos/{id:[0-9]+}",
		handleDeletePhoto,
	},
	Route{
		"users",
		[]string{"POST"},
		"/api/users",
		handlePostOrPutUser,
	},
	Route{
		"users",
		[]string{"POST"},
		"/api/users",
		handlePostOrPutUser,
	},
	Route{
		"users",
		[]string{"PUT"},
		"/api/users/{id:[0-9]+}",
		handlePostOrPutUser,
	},
	Route{
		"users",
		[]string{"DELETE"},
		"/api/users/{id:[0-9]+}",
		handleDeleteUser,
	},
	Route{
		"albums",
		[]string{"GET", "HEAD"},
		"/api/albums",
		handleGetAlbums,
	},
	Route{
		"albums",
		[]string{"POST"},
		"/api/albums",
		handlePostAlbums,
	},
	Route{
		"albums",
		[]string{"PUT"},
		"/api/albums/{name:[A-z][^/]*}",
		handlePutAlbums,
	},

	Route{
		"albums",
		[]string{"GET"},
		"/api/albums/{album_name:[A-z][^/]*}",
		handleGetAlbum,
	},
}
