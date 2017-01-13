// coalesce/config.go

package main

var cfg Config

type Config struct {
	System struct {
		Database     string
		DatabaseInit string
		ErrorLog     string
		AccessLog    string
		ResourceDir  string
		Listen       string
	}
}
