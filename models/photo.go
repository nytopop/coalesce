// coalesce/models

// photo.go

package models

import (
	"time"

	"github.com/nytopop/utron"
	"gopkg.in/mgo.v2/bson"
)

// represents a photo meta file
type Photo struct {
	Id         bson.ObjectId `schema:"id" bson:"_id,omitempty"`
	Length     int           `bson:"length"`
	ChunkSize  int           `bson:"chunkSize"`
	UploadDate time.Time     `schema:"uploadDate" bson:"uploadDate"`
	Md5        string        `schema:"md5" bson:"md5"`
	Filename   string        `schema:"filename" bson:"filename"`
}

type PhotoData struct {
}

func init() {
	utron.RegisterModels(&Photo{})
}
