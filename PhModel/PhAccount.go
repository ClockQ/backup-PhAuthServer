package PhModel

import (
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	EmployeeID string `json:"employee-id" bson:"employee-id"`
	Email      string `json:"email" bson:"email"`
	Phone      string `json:"phone" bson:"phone"`
	Username   string `json:"username" bson:"username"`
	Password   string `json:"password" bson:"password"`

	RoleID string `json:"-" bson:"role-id"`
	Role   *Role  `json:"-"`

	RegisterDate float64 `json:"register-date" bson:"register-date"`
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

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Account) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "roles",
			Name: "role",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Account) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.RoleID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.RoleID,
			Type: "roles",
			Name: "role",
		})
	}

	return result
}

func (u Account) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.RoleID != "" && u.Role != nil {
		result = append(result, u.Role)
	}

	return result
}

func (a *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
