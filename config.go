// coalesce/config.go

package main

// config file
type Config struct {
	Database struct {
		Host     string
		Name     string
		Username string
		Password string
	}
	Server struct {
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
