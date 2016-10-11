# coalesce

[![Build Status](https://travis-ci.org/nytopop/coalesce.svg?branch=master)](https://travis-ci.org/nytopop/coalesce)

A lightning fast multi-user blogging platform in Go.

## Deployment

Coalesce runs as a single binary using the native net/http webserver. It is designed to run behind a reverse proxy such as haproxy / nginx / apache. A mongodb instance is required.

## Technologies

[Gin Web Framework](https://github.com/gin-gonic/gin)

MongoDB via [mgo](https://github.com/go-mgo/mgo)

Tag generation via [cortical.io](http://www.cortical.io/)

Markdown processing via [Blackfriday](https://github.com/russross/blackfriday)

HTML5 template from [HTML5 UP](https://html5up.net/)

Configuration with [gcfg](https://github.com/go-gcfg/gcfg/tree/v1.2.0)
