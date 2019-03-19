package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Nickname string `json:"nickname" bson:"nickname"`
	Phone    string `json:"phone" bson:"phone"`

	Token string `json:"token"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (a Account) GetID() string {
	return a.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (a *Account) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
