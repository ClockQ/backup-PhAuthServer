package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	ID           string        `json:"-"`
	Id_          bson.ObjectId `json:"-" bson:"_id"`
	CompanyID    string        `json:"company-id" bson:"company-id"`
	Type         string        `json:"type" bson:"type"`
	Name         string        `json:"name" bson:"name"`
	Describe     string        `json:"describe" bson:"describe"`
	Parental     string        `json:"parental" bson:"parental"`
	RegisterDate float64        `json:"register-date" bson:"register-date"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Group) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Group) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Group) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
