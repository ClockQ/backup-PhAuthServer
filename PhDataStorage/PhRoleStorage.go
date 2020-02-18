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

// PhRoleStorage stores all of the tasty modelleaf, needs to be injected into
// Role and Role Resource. In the real world, you would use a database for that.
type PhRoleStorage struct {
	db *BmMongodb.BmMongodb
}

func (s PhRoleStorage) NewStorage(args []BmDaemons.BmDaemon) *PhRoleStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &PhRoleStorage{db: mdb}
}

// GetAll of the modelleaf
func (s PhRoleStorage) GetAll(r api2go.Request, skip int, take int) []PhModel.Role {
	in := PhModel.Role{}
	var out []PhModel.Role
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
func (s PhRoleStorage) GetOne(id string) (PhModel.Role, error) {
	in := PhModel.Role{ID: id}
	out := PhModel.Role{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Role for id %s not found", id)
	return PhModel.Role{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *PhRoleStorage) Insert(c PhModel.Role) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *PhRoleStorage) Delete(id string) error {
	in := PhModel.Role{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Role with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *PhRoleStorage) Update(c PhModel.Role) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Role with id does not exist")
	}

	return nil
}
