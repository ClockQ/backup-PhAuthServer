package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Role struct {
	ID           string        `json:"-"`
	Id_          bson.ObjectId `json:"-" bson:"_id"`
	Title        string        `json:"title" bson:"title"`
	Level        float64        `json:"level" bson:"level"`
	RegisterDate float64       `json:"register-date" bson:"register-date"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Role) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Role) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Role) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
