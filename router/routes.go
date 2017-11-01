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
	methods []string
	pattern string
	handler http.HandlerFunc
}

// Routes is a list of routes.
type Routes []Route

// Define the actual routes here.
var routes = Routes{
	Route{
		[]string{"GET"},
		"/auth/google/login",
		oauthRedirect,
	},
	Route{
		[]string{"GET"},
		"/auth/google/return",
		handleOauthCallback,
	},
	Route{
		[]string{"GET"},
		"/auth/logout",
		handleLogout,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/photos/{id:[0-9]+}",
		handlePhoto,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/photos/{id:[0-9]+}/image",
		handlePhotoImage,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/photos/{id:[0-9]+}/thumbnail",
		handlePhotoThumbnail,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/photos/{id:[0-9]+}/big_thumbnail",
		handlePhotoBigThumbnail,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/photos",
		handlePhotos,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/users",
		handleUsers,
	},
	Route{
		[]string{"POST"},
		"/photos",
		handlePostPhoto,
	},
	Route{
		[]string{"PUT"},
		"/photos/{id:[0-9]+}",
		handlePutPhoto,
	},
	Route{
		[]string{"DELETE"},
		"/photos/{id:[0-9]+}",
		handleDeletePhoto,
	},
	Route{
		[]string{"POST"},
		"/users",
		handlePostOrPutUser,
	},
	Route{
		[]string{"POST"},
		"/users",
		handlePostOrPutUser,
	},
	Route{
		[]string{"PUT"},
		"/users/{id:[0-9]+}",
		handlePostOrPutUser,
	},
	Route{
		[]string{"DELETE"},
		"/users/{id:[0-9]+}",
		handleDeleteUser,
	},
	Route{
		[]string{"GET", "HEAD"},
		"/albums",
		handleGetAlbums,
	},
	Route{
		[]string{"POST"},
		"/albums",
		handlePostAlbums,
	},
	Route{
		[]string{"PUT"},
		"/albums/{name:[A-z][^/]*}",
		handlePutAlbums,
	},

	Route{
		[]string{"GET"},
		"/albums/{album_name:[A-z][^/]*}",
		handleGetAlbum,
	},
}
