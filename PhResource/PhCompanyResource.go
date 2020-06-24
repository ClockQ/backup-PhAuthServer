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

type PhCompanyResource struct {
	PhGroupStorage *PhDataStorage.PhGroupStorage
	PhCompanyStorage *PhDataStorage.PhCompanyStorage
}

func (c PhCompanyResource) NewResource(args []BmDataStorage.BmStorage) *PhCompanyResource {
	var cs *PhDataStorage.PhCompanyStorage
	var gs *PhDataStorage.PhGroupStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhCompanyStorage" {
			cs = arg.(*PhDataStorage.PhCompanyStorage)
		} else if tp.Name() == "PhGroupStorage" {
			gs = arg.(*PhDataStorage.PhGroupStorage)
		}
	}
	return &PhCompanyResource{
		PhCompanyStorage: cs,
		PhGroupStorage: gs,
	}
}

// FindAll images
func (c PhCompanyResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	groupsID, ok := r.QueryParams["groupsID"]
	if ok {
		modelRootID := groupsID[0]
		modelRoot, err := c.PhGroupStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.CompanyID
		if modelID != "" {
			model, err := c.PhCompanyStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	var result []PhModel.Company
	result = c.PhCompanyStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhCompanyResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.PhCompanyStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new account
func (c PhCompanyResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Company)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.PhCompanyStorage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhCompanyResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhCompanyStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a account
func (c PhCompanyResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(PhModel.Company)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.PhCompanyStorage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
