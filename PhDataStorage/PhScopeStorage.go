package PhDataStorage

import (
	"errors"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/manyminds/api2go"
	"net/http"
	"ph_auth/PhModel"
)

// PhScopeStorage stores all of the tasty modelleaf, needs to be injected into
// Scope and Scope Resource. In the real world, you would use a database for that.
type PhScopeStorage struct {
	db *BmMongodb.BmMongodb
}

func (s PhScopeStorage) NewAccountStorage(args []BmDaemons.BmDaemon) *PhScopeStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &PhScopeStorage{db: mdb}
}

// GetAll of the modelleaf
func (s PhScopeStorage) GetAll(r api2go.Request, skip int, take int) []PhModel.Scope {
	in := PhModel.Scope{}
	var out []PhModel.Scope
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		for i, iter := range out {
			s.db.ResetIdWithId_(&iter)
			out[i] = iter
		}
		return out
	} else {
		return nil
	}
}

// GetOne tasty modelleaf
func (s PhScopeStorage) GetOne(id string) (PhModel.Scope, error) {
	in := PhModel.Scope{ID: id}
	out := PhModel.Scope{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Scope for id %s not found", id)
	return PhModel.Scope{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *PhScopeStorage) Insert(c PhModel.Scope) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *PhScopeStorage) Delete(id string) error {
	in := PhModel.Scope{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Scope with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *PhScopeStorage) Update(c PhModel.Scope) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Scope with id does not exist")
	}

	return nil
}
