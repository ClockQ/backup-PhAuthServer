package PhModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Company struct {
	ID  		string        `json:"-"`
	Id_ 		bson.ObjectId `json:"-" bson:"_id"`
	Name		string	`json:"name" bson:"name"`
	Level		string	`json:"level" bson:"level"`
	Location	string	`json:"location" bson:"location"`
	Industry	string	`json:"industry" bson:"industry"`
	Describe	string	`json:"describe" bson:"describe"`
	ImageID		string	`json:"-" bson:"image-id"`
	Image		*Image 	`json:"-"`
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

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Company) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "images",
			Name: "image",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Company) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.ImageID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.ImageID,
			Type: "images",
			Name: "image",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Company) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.ImageID != "" && u.Image != nil {
		result = append(result, u.Image)
	}

	return result
}

func (u *Company) SetToOneReferenceID(name, ID string) error {
	if name == "image" {
		u.ImageID = ID
		return nil
	}
	return errors.New("There is no to-one relationship with the name " + name)
}