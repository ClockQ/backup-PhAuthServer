package PhModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	EmployeeID string    `json:"-" bson:"employee-id"`
	Employee   *Employee `json:"-"`
	Email      string    `json:"-" bson:"email"`
	Phone      string    `json:"-" bson:"phone"`
	Username   string    `json:"username" bson:"username"`
	Password   string    `json:"-" bson:"password"`

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
		{
			Type: "employees",
			Name: "employee",
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

	if u.EmployeeID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.EmployeeID,
			Type: "employees",
			Name: "employee",
		})
	}

	return result
}

func (u Account) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.RoleID != "" && u.Role != nil {
		result = append(result, u.Role)
	}

	if u.EmployeeID != "" && u.Employee != nil {
		result = append(result, u.Employee)
	}

	return result
}

func (u *Account) SetToOneReferenceID(name, ID string) error {
	if name == "employee" {
		u.EmployeeID = ID
		return nil
	}
	if name == "role" {
		u.RoleID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

func (a *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	r := make(map[string]interface{})
	var ids []bson.ObjectId
	for k, v := range parameters {
		switch k {
		case "ids":
			for i := 0; i < len(v); i++ {
				ids = append(ids, bson.ObjectIdHex(v[i]))
			}
			r["$in"] = ids
			rst["_id"] = r
		}
	}
	return rst
}
