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

// PhEmployeeStroage stores all of the tasty modelleaf, needs to be injected into
// Employee and Employee Resource. In the real world, you would use a database for that.
type PhEmployeeStroage struct {
	db *BmMongodb.BmMongodb
}

func (s PhEmployeeStroage) NewStorage(args []BmDaemons.BmDaemon) *PhEmployeeStroage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &PhEmployeeStroage{db: mdb}
}

// GetAll of the modelleaf
func (s PhEmployeeStroage) GetAll(r api2go.Request, skip int, take int) []PhModel.Employee {
	in := PhModel.Employee{}
	var out []PhModel.Employee
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
func (s PhEmployeeStroage) GetOne(id string) (PhModel.Employee, error) {
	in := PhModel.Employee{ID: id}
	out := PhModel.Employee{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Employee for id %s not found", id)
	return PhModel.Employee{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *PhEmployeeStroage) Insert(c PhModel.Employee) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *PhEmployeeStroage) Delete(id string) error {
	in := PhModel.Employee{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Employee with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *PhEmployeeStroage) Update(c PhModel.Employee) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Employee with id does not exist")
	}

	return nil
}
