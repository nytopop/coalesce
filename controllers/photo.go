// coalesce/controllers

// photo.go

package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"git.echoesofthe.net/nytopop/coalesce/models"
	"github.com/nytopop/utron"
)

type PhotoController struct {
	*utron.BaseController
	Routes []string
}

// list of all photos
// change the template so it is a grid layout of latest 25 photos
func (p *PhotoController) Home() {
	// collection name needs to be user.files
	gfs := p.Ctx.DB.DB.GridFS("fs")

	photos := []*models.Photo{}
	err := gfs.Find(nil).Sort("uploadDate").Iter().All(&photos)
	if err != nil {
		log.Println("Error in database query")
	}

	p.Ctx.Data["List"] = photos
	p.Ctx.Template = "photo"
	p.HTML(http.StatusOK)
}

func (p *PhotoController) Create() {
	//gfs := p.Ctx.DB.DB.GridFS("fs")
	//decoder := schema.NewDecoder()

	req := p.Ctx.Request()
	//err :=
	if err := req.ParseMultipartForm(16777216); err != nil {
		fmt.Println(err.Error())
	}

	/*
		photo, err := gfs.Create(p.Ctx.Params["filename"])
		if err != nil {
			p.Ctx.Data["Message"] = err.Error()
			p.Ctx.Template = "error"
			p.HTML(http.StatusInternalServerError)
			return
		}
	*/

	// iterate through all file maps in the post form
	// this can be used for bulk uploads!
	for _, v := range req.MultipartForm.File {
		// iterate through each file map
		// empty []byte
		buf := new(bytes.Buffer)

		// open File object to file
		file, _ := v[0].Open()

		// attempt to read file into memory
		buf.ReadFrom(file)
		data := buf.Bytes()
		fmt.Println(v[0].Filename, len(data))
	}

	// req.MultipartForm.Value
	// contains form data

	// use gridfs to upload file
	// get post data
	// put into correct format
	// collect metadata
	// upload to database
	p.Ctx.Redirect("/photo", http.StatusFound)
}

func (p *PhotoController) Delete() {
	// open the currently active user's collection
	gfs := p.Ctx.DB.DB.GridFS("fs")

	id := p.Ctx.Params["id"]

	// get objectid obj from hex string
	bid := bson.ObjectIdHex(id)

	if err := gfs.RemoveId(bid); err != nil {
		p.Ctx.Data["Message"] = err.Error()
		p.Ctx.Template = "error"
		p.HTML(http.StatusInternalServerError)
		return
	}

	p.Ctx.Redirect("/photo", http.StatusFound)
}

// update routes so that del and upl are both post
func NewPhotoController() *PhotoController {
	return &PhotoController{
		Routes: []string{
			"get;/photo;Home",
			"post;/photo/new;Create",
			"get;/photo/del/{id};Delete",
		},
	}
}

func init() {
	utron.RegisterController(NewPhotoController())
}
