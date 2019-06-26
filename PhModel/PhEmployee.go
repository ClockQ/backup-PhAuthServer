package PhModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Employee struct {
	ID      string        `json:"-"`
	Id_     bson.ObjectId `json:"-" bson:"_id"`
	GroupID string        `json:"group-id" bson:"group-id"`
	Group   *Group        `json:"-"`
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

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Employee) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "groups",
			Name: "group",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Employee) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.GroupID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.GroupID,
			Type: "groups",
			Name: "group",
		})
	}

	return result
}

func (u Employee) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.GroupID != "" && u.Group != nil {
		result = append(result, u.Group)
	}

	return result
}

func (u *Employee) SetToOneReferenceID(name, ID string) error {
	if name == "group" {
		u.GroupID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

func (a *Employee) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
