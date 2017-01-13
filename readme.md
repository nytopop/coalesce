# coalesce [![Build Status](https://travis-ci.org/nytopop/coalesce.svg?branch=master)](https://travis-ci.org/nytopop/coalesce)

A lightning fast multi-user blogging platform in Go. Currently in pre-alpha phase.

## Installation

Coalesce runs as a single binary using the native net/http webserver. It is designed to run behind a reverse proxy such as haproxy / nginx / apache, but can be used standalone.

To build the statically linked .deb package (amd64):

    git clone https://github.com/nytopop/coalesce
    cd coalesce
    make

If make succeeded, the `build/` directory should contain a `dpkg -i` installable .deb package.

An administrator user with the credentials `admin` : `coalesce` is created by default on first run, *change the password immediately*. Any additional users that are created by registering will be initialized as `commentors`, which can be promoted to `editors` through the administrator account. There can only be one administrator (for now).

To run without installing:

    git clone https://github.com/nytopop/coalesce
    cd coalesce
    make run

## TODO

- [x] Comment delete
- [x] User delete
- [x] Post delete with comments
- [ ] Image CRUD
- [ ] Search
- [x] Display author / timestamp of posts in list and view
- [x] Cascading delete
- [ ] /users/me page, # comments, # posts, user information
- [ ] /users/mycomments
- [ ] Markdown comments
- [x] create admin user on initial load
- [ ] backup / restore database
- [ ] Next / Prev buttons on /posts and others
- [ ] Configuration of session secret
- [ ] Mailer
- [ ] OAuth 2.0 support for sign in
- [ ] RSS Feeds @ /rss
- [x] User password change support @ /users/me
- [ ] Site statistics @ /stats

## Packaging

- [x] .deb package build in makefile
- [x] binary location   (/usr/bin/coalesce)
- [x] config location   (/etc/coalesce.conf) -cfg flag
- [x] logging location  (/var/log/coalesce/err.log|access.log)
- [x] database location (/var/lib/coalesce.db)
- [x] resource location (/usr/share/coalesce)
- [x] systemd / sysv service files
- [ ] properly set file ownership for all resources

## Deps

Docker

Web framework [gin-gonic/gin](https://github.com/gin-gonic/gin)

SQLite3 via [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

Markdown processing via [russross/blackfriday](https://github.com/russross/blackfriday)

Image thumbnailing via [nfnt/resize](https://github.com/nfnt/resize)
