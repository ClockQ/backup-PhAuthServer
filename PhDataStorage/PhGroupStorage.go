package PhDataStorage

import (
	"errors"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/manyminds/api2go"
	"net/http"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
)

// PhGroupStorage stores all of the tasty modelleaf, needs to be injected into
// Group and Group Resource. In the real world, you would use a database for that.
type PhGroupStorage struct {
	db *BmMongodb.BmMongodb
}

func (s PhGroupStorage) NewStorage(args []BmDaemons.BmDaemon) *PhGroupStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &PhGroupStorage{db: mdb}
}

// GetAll of the modelleaf
func (s PhGroupStorage) GetAll(r api2go.Request, skip int, take int) []PhModel.Group {
	in := PhModel.Group{}
	var out []PhModel.Group
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
func (s PhGroupStorage) GetOne(id string) (PhModel.Group, error) {
	in := PhModel.Group{ID: id}
	out := PhModel.Group{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Group for id %s not found", id)
	return PhModel.Group{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *PhGroupStorage) Insert(c PhModel.Group) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *PhGroupStorage) Delete(id string) error {
	in := PhModel.Group{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Group with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *PhGroupStorage) Update(c PhModel.Group) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Group with id does not exist")
	}

	return nil
}
