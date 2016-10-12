// coalesce/config.go

package main

// config file
type Config struct {
	Server struct {
		ApiKey        string
		AdminPassword string
		Hostname      string
		Static        string
		Template      string
	}
	Site struct {
		Title       string
		Description string
		Owner       string
		Github      string
		Email       string
	}
}
