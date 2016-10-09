// coalesce/img.go

package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type Img struct {
	Id         bson.ObjectId `schema:"id" bson:"_id,omitempty"`
	Length     int           `bson:"length"`
	ChunkSize  int           `bson:"chunkSize"`
	UploadDate time.Time     `schema:"uploadDate" bson:"uploadDate"`
	Md5        string        `schema:"md5" bson:"md5"`
	Filename   string        `schema:"filename" bson:"filename"`
}

// home
func ImgHome(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(cfg.Database.Name)
	gfs := db.GridFS("fs")

	imgs := []*Img{}
	if err := gfs.Find(nil).Sort("uploadDate").Iter().All(&imgs); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(http.StatusOK, "img/home.html", gin.H{
		"Site": cfg.Site,
		"List": imgs,
		"User": GetUser(c),
	})
}

// thumb
func ImgThumb(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(cfg.Database.Name)
	gfs := db.GridFS("fs")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	// open photo by id
	file, err := gfs.OpenId(id)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// make image
	img, _, err := image.Decode(file)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}
	file.Close()

	// resize
	m := resize.Resize(500, 0, img, resize.Lanczos3)

	// encode
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// send
	c.Data(http.StatusOK, "image/jpg", buf.Bytes())
}

// view
func ImgView(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(cfg.Database.Name)
	gfs := db.GridFS("fs")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	// open photo by id
	file, err := gfs.OpenId(id)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// format
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	file.Close()

	// send
	c.Data(http.StatusOK, "image/jpg", buf.Bytes())
}

// upload
func ImgUpload(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(cfg.Database.Name)
	gfs := db.GridFS("fs")

	//req := c.Request
	if err := c.Request.ParseMultipartForm(36777216); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	for _, v := range c.Request.MultipartForm.File {
		for _, e := range v {
			buf := new(bytes.Buffer)

			// do batch processing
			file, err := e.Open()
			if err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			// read into memory
			if _, err := buf.ReadFrom(file); err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			// create empty file in gfs
			img, err := gfs.Create(e.Filename)
			if err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			// write to db file
			if _, err := img.Write(buf.Bytes()); err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			// close
			if err := img.Close(); err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}
		}
	}

	c.Redirect(302, "/img")
}
