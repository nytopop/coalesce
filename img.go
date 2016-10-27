// coalesce/img.go

package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type Img struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Length     int           `bson:"length"`
	ChunkSize  int           `bson:"chunkSize"`
	UploadDate time.Time     `bson:"uploadDate"`
	Md5        string        `bson:"md5"`
	Filename   string        `bson:"filename"`
}

// GET /img
func ImgAll(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(dbname)
	gfs := db.GridFS("images")

	imgs := []*Img{}
	if err := gfs.Find(nil).Sort("-uploadDate").Iter().All(&imgs); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(http.StatusOK, "img/all.html", gin.H{
		"Site": GetConf(),
		"User": GetUser(c),
		"Imgs": imgs,
	})
}

// GET /img/thumb/:id
func ImgThumb(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(dbname)
	gfs := db.GridFS("thumbs")

	hexid := c.Param("id")

	// open photo by id
	file, err := gfs.Open(hexid)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}
	defer file.Close()

	// format
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// send
	c.Data(http.StatusOK, "image/jpg", buf.Bytes())
}

// GET /img/view/:id
func ImgView(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(dbname)
	gfs := db.GridFS("images")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	// open img by id
	file, err := gfs.OpenId(id)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}
	defer file.Close()

	// format
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	// send
	c.Data(http.StatusOK, "image/jpg", buf.Bytes())
}

// GET /img/del/:id
func ImgTryDelete(c *gin.Context) {
	session := globalSession.Copy()
	db := session.DB(dbname)
	imgfs := db.GridFS("images")
	thumbfs := db.GridFS("thumbs")

	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	if err := imgfs.RemoveId(id); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	if err := thumbfs.Remove(hexid); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.Redirect(302, "/img")
}

// GET /img/new
func ImgNew(c *gin.Context) {
	c.HTML(http.StatusOK, "img/new.html", gin.H{
		"Site": GetConf(),
		"User": GetUser(c),
	})
}

// POST /img/new
func ImgTryNew(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(36777216); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	for _, v := range c.Request.MultipartForm.File {
		for _, e := range v {
			err, id := WriteImage(e)
			if err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			if err := WriteThumb(e, id); err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}
		}
	}

	c.Redirect(302, "/img")
}

func WriteImage(h *multipart.FileHeader) (error, string) {
	session := globalSession.Copy()
	db := session.DB(dbname)
	gfs := db.GridFS("images")

	file, err := h.Open()
	if err != nil {
		return err, ""
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err, ""
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return err, ""
	}

	i, err := gfs.Create(h.Filename)
	if err != nil {
		return err, ""
	}

	if _, err := i.Write(buf.Bytes()); err != nil {
		return err, ""
	}

	if err := i.Close(); err != nil {
		return err, ""
	}

	id := i.Id().(bson.ObjectId).Hex()

	return nil, id
}

func WriteThumb(h *multipart.FileHeader, id string) error {
	session := globalSession.Copy()
	db := session.DB(dbname)
	gfs := db.GridFS("thumbs")

	file, err := h.Open()
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	m := resize.Resize(500, 0, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		return err
	}

	i, err := gfs.Create(id)
	if err != nil {
		return err
	}

	if _, err := i.Write(buf.Bytes()); err != nil {
		return err
	}

	if err := i.Close(); err != nil {
		return err
	}

	return nil
}
