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

type PhGroupResource struct {
	PhEmployeeStroage *PhDataStorage.PhEmployeeStroage
	PhGroupStroage *PhDataStorage.PhGroupStroage
	PhCompanyStroage *PhDataStorage.PhCompanyStroage
}

func (c PhGroupResource) NewResource(args []BmDataStorage.BmStorage) *PhGroupResource {
	var es *PhDataStorage.PhEmployeeStroage
	var gs *PhDataStorage.PhGroupStroage
	var cs *PhDataStorage.PhCompanyStroage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhGroupStroage" {
			gs = arg.(*PhDataStorage.PhGroupStroage)
		} else if tp.Name() == "PhCompanyStroage" {
			cs = arg.(*PhDataStorage.PhCompanyStroage)
		} else if tp.Name() == "PhEmployeeStroage" {
			es = arg.(*PhDataStorage.PhEmployeeStroage)
		}
	}
	return &PhGroupResource{
		PhGroupStroage: gs,
		PhCompanyStroage: cs,
		PhEmployeeStroage: es,
	}
}

// FindAll images
func (c PhGroupResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	employeeId, ok := r.QueryParams["employeesID"]
	if ok {
		modelRootID := employeeId[0]
		modelRoot, err := c.PhEmployeeStroage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.GroupID
		if modelID != "" {
			model, err := c.PhGroupStroage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	var result []PhModel.Group
	result = c.PhGroupStroage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhGroupResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := c.PhGroupStroage.GetOne(ID)

	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.CompanyID != "" {
		company, err := c.PhCompanyStroage.GetOne(model.CompanyID)
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

	id := c.PhGroupStroage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhGroupResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhGroupStroage.Delete(id)
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

	err := c.PhGroupStroage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
