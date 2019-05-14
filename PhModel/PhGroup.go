package PhModel

import (
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	ID           string        `json:"-"`
	Id_          bson.ObjectId `json:"-" bson:"_id"`
	CompanyID    string        `json:"company-id" bson:"company-id"`
	Company      *Company      `json:"-"`
	Type         string        `json:"type" bson:"type"`
	Name         string        `json:"name" bson:"name"`
	Describe     string        `json:"describe" bson:"describe"`
	Parental     string        `json:"parental" bson:"parental"`
	RegisterDate float64       `json:"register-date" bson:"register-date"`
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

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Group) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "companies",
			Name: "company",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Group) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.CompanyID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.CompanyID,
			Type: "companies",
			Name: "company",
		})
	}

	return result
}

func (u Group) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.CompanyID != "" && u.Company != nil {
		result = append(result, u.Company)
	}

	return result
}

func (a *Group) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
