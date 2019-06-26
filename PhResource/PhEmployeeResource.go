package PhResource

import (
	"errors"
	"github.com/alfredyang1986/BmServiceDef/BmDataStorage"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"

	"ph_auth/PhDataStorage"
	"ph_auth/PhModel"
)

type PhEmployeeResource struct {
	PhEmployeeStroage *PhDataStorage.PhEmployeeStroage
	PhGroupStroage *PhDataStorage.PhGroupStroage
}

func (c PhEmployeeResource) NewResource(args []BmDataStorage.BmStorage) *PhEmployeeResource {
	var cs *PhDataStorage.PhEmployeeStroage
	var gs *PhDataStorage.PhGroupStroage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhEmployeeStroage" {
			cs = arg.(*PhDataStorage.PhEmployeeStroage)
		} else if tp.Name() == "PhGroupStroage" {
			gs = arg.(*PhDataStorage.PhGroupStroage)
		}
	}
	return &PhEmployeeResource{
		PhEmployeeStroage: cs,
		PhGroupStroage: gs,
	}
}

// FindAll images
func (c PhEmployeeResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []PhModel.Employee
	result = c.PhEmployeeStroage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhEmployeeResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := c.PhEmployeeStroage.GetOne(ID)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.GroupID != "" {
		g, err := c.PhGroupStroage.GetOne(model.GroupID)
		if err != nil {
			return &Response{}, err
		}
		model.Group = &g
	}

	return &Response{Res: model}, err
}

// Create a new account
func (c PhEmployeeResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Employee)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.PhEmployeeStroage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhEmployeeResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhEmployeeStroage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a account
func (c PhEmployeeResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Employee)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.PhEmployeeStroage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
