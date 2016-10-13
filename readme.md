# coalesce

[![Build Status](https://travis-ci.org/nytopop/coalesce.svg?branch=master)](https://travis-ci.org/nytopop/coalesce)

A lightning fast multi-user blogging platform in Go.

## Installation

Coalesce runs as a single binary using the native net/http webserver. It is designed to run behind a reverse proxy such as haproxy / nginx / apache. A mongodb instance is required.

To run coalesce/mongodb using docker:

    git clone https://github.com/nytopop/coalesce
    cd coalesce
    make
    make run

Coalesce will be available on port :8080

## Technologies

Docker

Web framework [gin-gonic/gin](https://github.com/gin-gonic/gin)

MongoDB via [go-mgo/mgo](https://github.com/go-mgo/mgo)

Tag generation via [cortical.io](http://www.cortical.io/)

Markdown processing via [russross/blackfriday](https://github.com/russross/blackfriday)

Image thumbnailing via [nfnt/resize](https://github.com/nfnt/resize)

Configuration via [go-gcfg/gcfg](https://github.com/go-gcfg/gcfg/tree/v1.2.0)

HTML5 template from [HTML5 UP](https://html5up.net/)
