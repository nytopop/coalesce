# Start from a Debian image with the latest version of Go installed
FROM golang:onbuild

# Document that the service listens on port 8080.
EXPOSE 8080
