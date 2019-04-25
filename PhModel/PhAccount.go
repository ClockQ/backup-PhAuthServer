package PhModel

import (
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Nickname string `json:"nickname" bson:"nickname"`
	Gender   int    `json:"gender" bson:"gender"`

	CompanyID string   `json:"-" bson:"company-id"`
	Company   *Company `json:"-"`

	DepartmentID string      `json:"-" bson:"department-id"`
	Department   *Department `json:"-"`

	//ClientID		string	`json:"-" bson:"client-id"`
	//Client			*Client	`json:"-"`

	ScopeIDs []string `json:"-" bson:"scope-ids"`
	Scopes   []*Scope `json:"-"`

	ImageID string `json:"-" bson:"image-id"`
	Image   *Image `json:"-"`

	Phone string `json:"phone" bson:"phone"`
	Scope string `json:"scope" bson:"scope"`
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
			Type: "images",
			Name: "image",
		},
		{
			Type: "companies",
			Name: "company",
		},
		{
			Type: "departments",
			Name: "department",
		},
		{
			Type: "clients",
			Name: "client",
		},
		{
			Type: "scopes",
			Name: "scope",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Account) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.ImageID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.ImageID,
			Type: "images",
			Name: "image",
		})
	}

	if u.CompanyID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.CompanyID,
			Type: "companies",
			Name: "company",
		})
	}

	if u.DepartmentID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.DepartmentID,
			Type: "departments",
			Name: "department",
		})
	}

	//if u.ClientID != "" {
	//	result = append(result, jsonapi.ReferenceID{
	//		ID:   u.ClientID,
	//		Type: "clients",
	//		Name: "client",
	//	})
	//}

	return result
}

func (u Account) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.ImageID != "" && u.Image != nil {
		result = append(result, u.Image)
	}

	if u.CompanyID != "" && u.Company != nil {
		result = append(result, u.Company)
	}

	if u.DepartmentID != "" && u.Department != nil {
		result = append(result, u.Department)
	}

	//if u.ClientID != "" && u.Client != nil {
	//	result = append(result, u.Client)
	//}

	for key := range u.ScopeIDs {
		result = append(result, u.Scopes[key])
	}
	return result
}

func (a *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
