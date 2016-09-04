// coalesce/controllers

// img.go

package controllers

import (
	"bytes"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"git.echoesofthe.net/nytopop/coalesce/models"
	"github.com/nytopop/utron"
)

type ImgController struct {
	*utron.BaseController
	Routes []string
}

func sendErr(p *ImgController, err error) {
	if err != nil {
		p.Ctx.Data["Message"] = err.Error()
		p.Ctx.Template = "error"
		p.HTML(http.StatusInternalServerError)
	}
}

func (p *ImgController) Home() {
	// collection name needs to be user.files
	gfs := p.Ctx.DB.DB.GridFS("fs")

	photos := []*models.Photo{}
	if err := gfs.Find(nil).Sort("uploadDate").Iter().All(&photos); err != nil {
		sendErr(p, err)
		return
	}

	p.Ctx.Data["List"] = photos
	p.Ctx.Template = "img"
	p.HTML(http.StatusOK)
}

func (p *ImgController) Create() {
	gfs := p.Ctx.DB.DB.GridFS("fs")

	req := p.Ctx.Request()
	if err := req.ParseMultipartForm(16777216); err != nil {
		sendErr(p, err)
		return
	}

	for _, v := range req.MultipartForm.File {
		buf := new(bytes.Buffer)

		// open File object to file
		file, err := v[0].Open()
		if err != nil {
			sendErr(p, err)
			return
		}

		// attempt to read file into memory
		if _, err := buf.ReadFrom(file); err != nil {
			sendErr(p, err)
			return
		}

		// create empty file in gfs
		photo, err := gfs.Create(v[0].Filename)
		if err != nil {
			sendErr(p, err)
			return
		}

		// write to db file
		if _, err := photo.Write(buf.Bytes()); err != nil {
			sendErr(p, err)
			return
		}

		// close
		if err := photo.Close(); err != nil {
			sendErr(p, err)
			return
		}
	}

	p.Ctx.Redirect("/img", http.StatusFound)
}

func (p *ImgController) Delete() {
	// open the currently active user's collection
	gfs := p.Ctx.DB.DB.GridFS("fs")

	// get objectid obj from hex string
	hexid := p.Ctx.Params["id"]
	id := bson.ObjectIdHex(hexid)

	if err := gfs.RemoveId(id); err != nil {
		sendErr(p, err)
		return
	}

	p.Ctx.Redirect("/img", http.StatusFound)
}

// view a single image by id
func (p *ImgController) View() {
	gfs := p.Ctx.DB.DB.GridFS("fs")

	// get obj id from hex
	hexid := p.Ctx.Params["id"]
	id := bson.ObjectIdHex(hexid)

	// open photo by id
	photo, err := gfs.OpenId(id)
	if err != nil {
		sendErr(p, err)
	}

	// read file into []byte
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(photo); err != nil {
		sendErr(p, err)
	}

	// convert into image / whatever format
	p.Ctx.SetHeader("Content-Type", "image/jpeg")
	p.Ctx.SetHeader("Content-Length", strconv.Itoa(len(buf.Bytes())))
	if _, err := p.Ctx.Write(buf.Bytes()); err != nil {
		sendErr(p, err)
		return
	}
}

func NewImgController() *ImgController {
	return &ImgController{
		Routes: []string{
			"get;/img;Home",
			"post;/img/new;Create",
			"get;/img/del/{id};Delete",
			"get;/img/view/{id};View",
		},
	}
}

func init() {
	utron.RegisterController(NewImgController())
}
