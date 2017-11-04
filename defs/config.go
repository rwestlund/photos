// Copyright (c) 2017, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package defs

import "github.com/BurntSushi/toml"

// Config holds the configuration, which is read from a file at startup.
type Config struct {
	// ListenAddress should be in the format: [address]:port
	// It is passed to http.ListenAndServe().
	ListenAddress string

	DatabaseUserName string
	DatabaseName     string

	// OAuthClientID is the Google OAuth client id.
	OAuthClientID string

	// OAuthClientSecret is the Google OAuth client secret.
	OAuthClientSecret string

	// LocalHostName is the hostname in the URL where this server can be found.
	LocalHostName string

	// DefaultAdmins contains users who will be added as admins during resetdb.
	DefaultAdmins []string
}

// ReadConfigFile reads it the configuration.
func ReadConfigFile() (*Config, error) {
	var c = Config{}
	var _, err = toml.DecodeFile("config.toml", &c)
	return &c, err
}
