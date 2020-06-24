package PhDataStorage

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/manyminds/api2go"
)

type PhApplyuserStorage struct {
	db *BmMongodb.BmMongodb
}

func (s PhApplyuserStorage) NewApplyuserStorage(args []BmDaemons.BmDaemon) *PhApplyuserStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &PhApplyuserStorage{mdb}
}

// GetAll of the modelleaf
func (s PhApplyuserStorage) GetAll(r api2go.Request, skip int, take int) []*PhModel.Applyuser {
	in := PhModel.Applyuser{}
	var out []PhModel.Applyuser
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*PhModel.Applyuser
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil
	}
}

// GetOne tasty modelleaf
func (s PhApplyuserStorage) GetOne(id string) (PhModel.Applyuser, error) {
	in := PhModel.Applyuser{ID: id}
	out := PhModel.Applyuser{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Applyuser for id %s not found", id)
	return PhModel.Applyuser{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *PhApplyuserStorage) Insert(c PhModel.Applyuser) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *PhApplyuserStorage) Delete(id string) error {
	in := PhModel.Applyuser{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Applyuser with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *PhApplyuserStorage) Update(c PhModel.Applyuser) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Applyuser with id does not exist")
	}

	return nil
}

func (s *PhApplyuserStorage) Count(req api2go.Request, c PhModel.Applyuser) int {
	r, _ := s.db.Count(req, &c)
	return r
}