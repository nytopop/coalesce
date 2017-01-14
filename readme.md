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

## TODO for alpha release

### Features
- [ ] Search   [/search/texttosearch]
- [ ] RSS Feed [/rss]
- [ ] Image CRUD [/img/new|del|view|thumb]
- [ ] Markdown for comments
- [ ] Backup / Restore / Reset database [/db/backup|restore|reset]
- [ ] Dedicated change password page, add to usermenu
- [ ] OAuth 2.0
- [ ] Mailer
- [ ] Site statistics [/stats]
- [ ] User profiles, name|email|about
- [ ] Automated spam prevention / detection

### Fixes
- [ ] Next / Prev buttons on /posts|/users/myposts|/users/mycomments
- [ ] Stable configuration of dynamic session secret
- [ ] Session expiry timeout

## Deps

Docker

Web framework [gin-gonic/gin](https://github.com/gin-gonic/gin)

SQLite3 via [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

Markdown processing via [russross/blackfriday](https://github.com/russross/blackfriday)

Image thumbnailing via [nfnt/resize](https://github.com/nfnt/resize)
