package PhResource

import (
	"errors"
	"github.com/alfredyang1986/BmServiceDef/BmDataStorage"
	"github.com/manyminds/api2go"
	"net/http"
	"github.com/PharbersDeveloper/PhAuthServer/PhDataStorage"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"reflect"
	"strconv"
)

type PhApplyuserResource struct {
	PhApplyuserStorage 	*PhDataStorage.PhApplyuserStorage
}

func (s PhApplyuserResource) NewApplyuserResource(args []BmDataStorage.BmStorage) *PhApplyuserResource {
	var bis *PhDataStorage.PhApplyuserStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "PhApplyuserStorage" {
			bis = arg.(*PhDataStorage.PhApplyuserStorage)
		}
	}
	return &PhApplyuserResource{
		PhApplyuserStorage: bis,
	}
}

func (s PhApplyuserResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []*PhModel.Applyuser

	result = s.PhApplyuserStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s PhApplyuserResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []PhModel.Applyuser
		number, size, offset, limit string
	)

	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}

	if size != "" {
		sizeI, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		numberI, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		start := sizeI * (numberI - 1)
		for _, iter := range s.PhApplyuserStorage.GetAll(r, int(start), int(sizeI)) {
			result = append(result, *iter)
		}

	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		for _, iter := range s.PhApplyuserStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := PhModel.Applyuser{}
	count := s.PhApplyuserStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s PhApplyuserResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.PhApplyuserStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s PhApplyuserResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(PhModel.Applyuser)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	//r.QueryParams["email"] = []string{model.Email}
	//result := s.PhApplyuserStorage.GetAll(r, -1,-1)
	//if len(result) > 0 {
	//	panic("邮箱已存在")
	//}

	id := s.PhApplyuserStorage.Insert(model)
	model.ID = id

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s PhApplyuserResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.PhApplyuserStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s PhApplyuserResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(PhModel.Applyuser)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.PhApplyuserStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}

