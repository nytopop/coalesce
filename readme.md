# coalesce [![Build Status](https://travis-ci.org/nytopop/coalesce.svg?branch=master)](https://travis-ci.org/nytopop/coalesce)

A lightning fast multi-user blogging platform in Go. Currently in pre-alpha phase.

## Installation

Coalesce runs as a single binary using the native net/http webserver. It is designed to run behind a reverse proxy such as haproxy / nginx / apache.

## TODO

- [ ] Comment delete
- [x] User delete
- [x] Post delete with comments
- [ ] Image CRUD
- [ ] Search
- [x] Display author / timestamp of posts in list and view
- [x] Cascading delete
- [ ] /users/me page, # comments, # posts, user information
- [ ] /users/mycomments
- [ ] Markdown commenting
- [ ] Clean up SQL interface
- [ ] systemd service files
- [ ] configuration file or store configuration in DB
- [ ] backup / restore database

## Technologies

Docker

Web framework [gin-gonic/gin](https://github.com/gin-gonic/gin)

SQLite3 via [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

Markdown processing via [russross/blackfriday](https://github.com/russross/blackfriday)

Image thumbnailing via [nfnt/resize](https://github.com/nfnt/resize)
