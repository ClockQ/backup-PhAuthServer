package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Company struct {
	ID           string                 `json:"-"`
	Id_          bson.ObjectId          `json:"-" bson:"_id"`
	Name         string                 `json:"name" bson:"name"`
	Shorten      string                 `json:"shorten" bson:"shorten"`
	Location     map[string]interface{} `json:"location" bson:"location"`
	Headquarters bool                   `json:"headquarters" bson:"headquarters"`
	Parental     string                 `json:"parental" bson:"parental"`
	Describe     string                 `json:"describe" bson:"describe"`
	Image        string                 `json:"image" bson:"image"`
	RegisterDate float64                 `json:"register-date" bson:"register-date"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Company) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Company) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Company) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
