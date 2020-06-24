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

// PhCompanyStorage stores all of the tasty modelleaf, needs to be injected into
// Company and Company Resource. In the real world, you would use a database for that.
type PhCompanyStorage struct {
	db *BmMongodb.BmMongodb
}

func (s PhCompanyStorage) NewStorage(args []BmDaemons.BmDaemon) *PhCompanyStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &PhCompanyStorage{db: mdb}
}

// GetAll of the modelleaf
func (s PhCompanyStorage) GetAll(r api2go.Request, skip int, take int) []PhModel.Company {
	in := PhModel.Company{}
	var out []PhModel.Company
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
func (s PhCompanyStorage) GetOne(id string) (PhModel.Company, error) {
	in := PhModel.Company{ID: id}
	out := PhModel.Company{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Company for id %s not found", id)
	return PhModel.Company{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *PhCompanyStorage) Insert(c PhModel.Company) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *PhCompanyStorage) Delete(id string) error {
	in := PhModel.Company{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Company with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *PhCompanyStorage) Update(c PhModel.Company) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Company with id does not exist")
	}

	return nil
}
