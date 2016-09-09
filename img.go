// coalesce/img.go

package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
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
		//handle the err
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "img/home.html", gin.H{
		"Site": cfg.Site,
		"List": imgs,
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
		log.Println(err)
		return
	}

	// make image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		return
	}
	file.Close()

	// resize
	m := resize.Resize(500, 0, img, resize.Lanczos3)

	// encode
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		log.Println(err)
		return
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
		log.Println(err)
		return
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
		log.Println(err)
		return
	}

	for _, v := range c.Request.MultipartForm.File {
		for _, e := range v {
			buf := new(bytes.Buffer)

			// do batch processing
			file, err := e.Open()
			if err != nil {
				log.Println(err)
				return
			}

			// read into memory
			if _, err := buf.ReadFrom(file); err != nil {
				log.Println(err)
				return
			}

			// create empty file in gfs
			img, err := gfs.Create(e.Filename)
			if err != nil {
				log.Println(err)
				return
			}

			// write to db file
			if _, err := img.Write(buf.Bytes()); err != nil {
				log.Println(err)
				return
			}

			// close
			if err := img.Close(); err != nil {
				log.Println(err)
				return
			}
		}
	}

	c.Redirect(302, "/img")
}
