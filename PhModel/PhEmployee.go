package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Employee struct {
	ID      string        `json:"-"`
	Id_     bson.ObjectId `json:"-" bson:"_id"`
	GroupID string        `json:"group-id" bson:"group-id"`
	Name    string        `json:"name" bson:"name"`
	Gender  float64       `json:"gender" bson:"gender"`
	Image   string        `json:"image" bson:"image"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Employee) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Employee) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Employee) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
