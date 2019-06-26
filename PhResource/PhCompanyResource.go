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

type PhCompanyResource struct {
	PhGroupStroage *PhDataStorage.PhGroupStroage
	PhCompanyStroage *PhDataStorage.PhCompanyStroage
}

func (c PhCompanyResource) NewResource(args []BmDataStorage.BmStorage) *PhCompanyResource {
	var cs *PhDataStorage.PhCompanyStroage
	var gs *PhDataStorage.PhGroupStroage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhCompanyStroage" {
			cs = arg.(*PhDataStorage.PhCompanyStroage)
		} else if tp.Name() == "PhGroupStroage" {
			gs = arg.(*PhDataStorage.PhGroupStroage)
		}
	}
	return &PhCompanyResource{
		PhCompanyStroage: cs,
		PhGroupStroage: gs,
	}
}

// FindAll images
func (c PhCompanyResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	groupsID, ok := r.QueryParams["groupsID"]
	if ok {
		modelRootID := groupsID[0]
		modelRoot, err := c.PhGroupStroage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.CompanyID
		if modelID != "" {
			model, err := c.PhCompanyStroage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	var result []PhModel.Company
	result = c.PhCompanyStroage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c PhCompanyResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.PhCompanyStroage.GetOne(ID)
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

	id := c.PhCompanyStroage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c PhCompanyResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.PhCompanyStroage.Delete(id)
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

	err := c.PhCompanyStroage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
