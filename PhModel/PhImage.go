package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Image struct {
	ID   string        `json:"-"`
	Id_  bson.ObjectId `json:"-" bson:"_id"`
	Img  string        `json:"img" bson:"img"`
	Tag  string        `json:"tag" bson:"tag"`
	Flag int           `json:"flag" bson:"flag"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Image) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Image) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Image) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
