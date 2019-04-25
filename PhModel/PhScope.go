package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Scope struct {
	ID       string        `json:"-"`
	Id_      bson.ObjectId `json:"-" bson:"_id"`
	Level    string        `json:"level" bson:"level"`
	Value    string        `json:"value" bson:"value"`
	Describe string        `json:"describe" bson:"describe"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Scope) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Scope) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Scope) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
