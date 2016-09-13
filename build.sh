#!/bin/sh

go fmt *.go
go build -o bin/coalesce && bin/coalesce
