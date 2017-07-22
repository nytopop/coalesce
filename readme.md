# coalesce [![Build Status](https://travis-ci.org/nytopop/coalesce.svg?branch=master)](https://travis-ci.org/nytopop/coalesce)

A lightning fast multi-user blogging platform in Go. Currently in pre-alpha phase.

## Installation

Coalesce runs as a single binary using the native net/http webserver. It is designed to run behind a reverse proxy such as haproxy / nginx / apache, but can be used standalone.

An administrator user with the credentials `admin` : `coalesce` is created by default on first run, *change the password asap*. Any additional users that are created by registering will be initialized as `commentors`, which can be promoted to `editors` through the administrator account. There can only be one administrator (for now).

## TODO for alpha release
### Features
- [ ] Search   [/search/texttosearch]
- [ ] RSS Feed [/rss]
- [ ] Backup / Restore / Reset database [/db/backup|restore|reset]
- [ ] Mailer
- [ ] Site statistics [/stats]
- [ ] User profiles, name|email|about
- [ ] Automated spam prevention / detection

### Fixes
- [ ] Next / Prev buttons on /posts|/users/myposts|/users/mycomments
- [ ] Stable configuration of dynamic session secret
- [ ] Session expiry timeout
- [ ] Dedicated change password page, add to usermenu

### Other
- [ ] build packages with FPM
- [ ] vendor dependencies

## Deps
Web framework [gin-gonic/gin](https://github.com/gin-gonic/gin)

SQLite3 via [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

Markdown processing via [russross/blackfriday](https://github.com/russross/blackfriday)

Image thumbnailing via [nfnt/resize](https://github.com/nfnt/resize)
