package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type OpenClient struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Secret string `json:"secret" bson:"secret"`
	Domain string `json:"domain" bson:"domain"`
	UserID string `json:"user-id" bson:"user-id"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a OpenClient) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *OpenClient) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *OpenClient) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
