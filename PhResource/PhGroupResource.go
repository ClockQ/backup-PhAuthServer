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

type PhGroupResource struct {
	PhEmployeeStorage *PhDataStorage.PhEmployeeStorage
	PhGroupStorage *PhDataStorage.PhGroupStorage
	PhCompanyStorage *PhDataStorage.PhCompanyStorage
}

func (c PhGroupResource) NewResource(args []BmDataStorage.BmStorage) *PhGroupResource {
	var es *PhDataStorage.PhEmployeeStorage
	var gs *PhDataStorage.PhGroupStorage
	var cs *PhDataStorage.PhCompanyStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhGroupStorage" {
			gs = arg.(*PhDataStorage.PhGroupStorage)
		} else if tp.Name() == "PhCompanyStorage" {
			cs = arg.(*PhDataStorage.PhCompanyStorage)
		} else if tp.Name() == "PhEmployeeStorage" {
			es = arg.(*PhDataStorage.PhEmployeeStorage)
		}
	}
	return &PhGroupResource{
		PhGroupStorage: gs,
		PhCompanyStorage: cs,
		PhEmployeeStorage: es,
	}
}

// FindAll images
func (c PhGroupResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	employeeId, ok := r.QueryParams["employeesID"]
	if ok {
		modelRootID := employeeId[0]
		modelRoot, err := c.PhEmployeeStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.GroupID
		if modelID != "" {
			model, err := c.PhGroupStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	var result []PhModel.Group
	result = c.PhGroupStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhGroupResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := c.PhGroupStorage.GetOne(ID)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.CompanyID != "" {
		company, err := c.PhCompanyStorage.GetOne(model.CompanyID)
		if err != nil {
			return &Response{}, err
		}
		model.Company = &company
	}

	return &Response{Res: model}, err
}

// Create a new account
func (c PhGroupResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Group)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.PhGroupStorage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhGroupResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhGroupStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a account
func (c PhGroupResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Group)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.PhGroupStorage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
