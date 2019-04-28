package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Scope struct {
	ID           string        `json:"-"`
	Id_          bson.ObjectId `json:"-" bson:"_id"`
	GroupID      string        `json:"group-id" bson:"group-id"`
	Access       string        `json:"access" bson:"access"`
	Operation    string        `json:"operation" bson:"operation"`
	Expired      float64       `json:"expired" bson:"expired"`
	RegisterDate float64       `json:"register-date" bson:"register-date"`
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
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "group-id":
			rst[k] = v[0]
		}
	}
	return rst
}
