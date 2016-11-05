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

/* Routes are a list of these structs. */
type Route struct {
	name    string
	methods []string
	pattern string
	handler http.HandlerFunc
}
type Routes []Route

/* Define the actual routes here. */
var routes = Routes{
	Route{
		"auth",
		[]string{"GET"},
		"/auth/google/login",
		oauth_redirect,
	},
	Route{
		"auth",
		[]string{"GET"},
		"/auth/google/return",
		handle_oauth_callback,
	},
	Route{
		"logout",
		[]string{"GET"},
		"/auth/logout",
		handle_logout,
	},
	Route{
		"photo",
		[]string{"GET", "HEAD"},
		"/api/photos/{id:[0-9]+}",
		handle_photo,
	},
	Route{
		"photo",
		[]string{"GET", "HEAD"},
		"/api/photos/{id:[0-9]+}/image",
		handle_photo_image,
	},

	Route{
		"photos",
		[]string{"GET", "HEAD"},
		"/api/photos",
		handle_photos,
	},
	Route{
		"users",
		[]string{"GET", "HEAD"},
		"/api/users",
		handle_users,
	},
	Route{
		"photos",
		[]string{"POST"},
		"/api/photos",
		handle_post_photo,
	},
	Route{
		"photos",
		[]string{"PUT"},
		"/api/photos/{id:[0-9]+}",
		handle_put_photo,
	},
	Route{
		"photos",
		[]string{"DELETE"},
		"/api/photos/{id:[0-9]+}",
		handle_delete_photo,
	},
	Route{
		"users",
		[]string{"POST"},
		"/api/users",
		handle_post_or_put_user,
	},
	Route{
		"users",
		[]string{"POST"},
		"/api/users",
		handle_post_or_put_user,
	},
	Route{
		"users",
		[]string{"PUT"},
		"/api/users/{id:[0-9]+}",
		handle_post_or_put_user,
	},
	Route{
		"users",
		[]string{"DELETE"},
		"/api/users/{id:[0-9]+}",
		handle_delete_user,
	},
	Route{
		"albums",
		[]string{"GET", "HEAD"},
		"/api/albums",
		handle_get_albums,
	},
	Route{
		"albums",
		[]string{"POST"},
		"/api/albums",
		handle_post_albums,
	},
	Route{
		"albums",
		[]string{"PUT"},
		"/api/albums/{name:[A-z][^/]*}",
		handle_put_albums,
	},

	Route{
		"albums",
		[]string{"GET"},
		"/api/albums/{album_name:[A-z][^/]*}",
		handle_get_album,
	},
}
