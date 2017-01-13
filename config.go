// coalesce/config.go

package main

var cfg Config

type Config struct {
	System struct {
		Database string
		Log      string
		Resource string
		Listen   string
	}
}
