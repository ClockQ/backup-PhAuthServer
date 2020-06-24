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

type PhRoleResource struct {
	PhRoleStorage *PhDataStorage.PhRoleStorage
}

func (c PhRoleResource) NewResource(args []BmDataStorage.BmStorage) *PhRoleResource {
	var cs *PhDataStorage.PhRoleStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhRoleStorage" {
			cs = arg.(*PhDataStorage.PhRoleStorage)
		}
	}
	return &PhRoleResource{
		PhRoleStorage: cs,
	}
}

// FindAll images
func (c PhRoleResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []PhModel.Role
	result = c.PhRoleStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhRoleResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.PhRoleStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new account
func (c PhRoleResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Role)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.PhRoleStorage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhRoleResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhRoleStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a account
func (c PhRoleResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Role)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.PhRoleStorage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
