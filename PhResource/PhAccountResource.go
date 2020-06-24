package PhResource

import (
	"errors"
	"github.com/alfredyang1986/BmServiceDef/BmDataStorage"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"

	"github.com/PharbersDeveloper/PhAuthServer/PhDataStorage"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
)

type PhAccountResource struct {
	PhAccountStorage *PhDataStorage.PhAccountStorage
	PhEmployeeStorage *PhDataStorage.PhEmployeeStorage
	PhRoleStorage *PhDataStorage.PhRoleStorage
}

func (c PhAccountResource) NewResource(args []BmDataStorage.BmStorage) *PhAccountResource {
	var cs *PhDataStorage.PhAccountStorage
	var es *PhDataStorage.PhEmployeeStorage
	var rs *PhDataStorage.PhRoleStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhAccountStorage" {
			cs = arg.(*PhDataStorage.PhAccountStorage)
		} else if tp.Name() == "PhEmployeeStorage" {
			es = arg.(*PhDataStorage.PhEmployeeStorage)
		} else if tp.Name() == "PhRoleStorage" {
			rs = arg.(*PhDataStorage.PhRoleStorage)
		}
	}
	return &PhAccountResource{
		PhAccountStorage: cs,
		PhEmployeeStorage: es,
		PhRoleStorage: rs,
	}
}

// FindAll images
func (c PhAccountResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []PhModel.Account
	result = c.PhAccountStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhAccountResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := c.PhAccountStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.EmployeeID != "" {
		e, err := c.PhEmployeeStorage.GetOne(model.EmployeeID)
		if err != nil {
			return &Response{}, err
		}
		model.Employee = &e
	}

	if model.RoleID != "" {
		role, err := c.PhRoleStorage.GetOne(model.RoleID)
		if err != nil {
			return &Response{}, err
		}
		model.Role = &role
	}

	return &Response{Res: model}, err
}

// Create a new account
func (c PhAccountResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Account)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.PhAccountStorage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhAccountResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhAccountStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a account
func (c PhAccountResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Account)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.PhAccountStorage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
