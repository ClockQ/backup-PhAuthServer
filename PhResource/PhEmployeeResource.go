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
	PhEmployeeStorage 	*PhDataStorage.PhEmployeeStorage
	PhGroupStorage 		*PhDataStorage.PhGroupStorage
	PhAccountStorage 	*PhDataStorage.PhAccountStorage
}

func (c PhEmployeeResource) NewResource(args []BmDataStorage.BmStorage) *PhEmployeeResource {
	var cs *PhDataStorage.PhEmployeeStorage
	var gs *PhDataStorage.PhGroupStorage
	var ac *PhDataStorage.PhAccountStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhEmployeeStorage" {
			cs = arg.(*PhDataStorage.PhEmployeeStorage)
		} else if tp.Name() == "PhGroupStorage" {
			gs = arg.(*PhDataStorage.PhGroupStorage)
		} else if tp.Name() == "PhAccountStorage" {
			ac = arg.(*PhDataStorage.PhAccountStorage)
		}
	}
	return &PhEmployeeResource{
		PhEmployeeStorage: cs,
		PhGroupStorage: gs,
		PhAccountStorage: ac,
	}
}

// FindAll Employees
func (c PhEmployeeResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	accountsId, ok := r.QueryParams["accountsID"]
	if ok {
		modelRootID := accountsId[0]
		modelRoot, err := c.PhAccountStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.EmployeeID
		model, err := c.PhEmployeeStorage.GetOne(modelID)
		if err != nil {
			return &Response{}, err
		}
		return &Response{Res: model}, nil
	}

	var result []PhModel.Employee
	result = c.PhEmployeeStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne Employee
func (c PhEmployeeResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := c.PhEmployeeStorage.GetOne(ID)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.GroupID != "" {
		g, err := c.PhGroupStorage.GetOne(model.GroupID)
		if err != nil {
			return &Response{}, err
		}
		model.Group = &g
	}

	return &Response{Res: model}, err
}

// Create a new Employee
func (c PhEmployeeResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Employee)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.PhEmployeeStorage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a Employee :(
func (c PhEmployeeResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhEmployeeStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a Employee
func (c PhEmployeeResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Employee)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.PhEmployeeStorage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
