package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Department struct {
	ID  		string        `json:"-"`
	Id_ 		bson.ObjectId `json:"-" bson:"_id"`
	Name 		string	`json:"name" bson:"name"`
	Describe	string	`json:"describe" bson:"describe"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Department) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Department) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Department) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
