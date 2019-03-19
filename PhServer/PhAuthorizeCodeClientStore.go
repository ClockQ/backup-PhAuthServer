package PhServer

import (
	"errors"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"github.com/alfredyang1986/BmServiceDef/BmModel"
	"github.com/PharbersDeveloper/PhAuthServer/PhModel"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	)

func NewAuthorizeCodeClientStore(mdb *BmMongodb.BmMongodb) *PhAuthorizeCodeClientStore {
	return &PhAuthorizeCodeClientStore{
		mdb: mdb,
	}
}

// PhAuthorizeCodeClientStore client information store
// PhAuthorizeCodeClientStore implement ==> oauth2.ClientInfo
type PhAuthorizeCodeClientStore struct {
	mdb *BmMongodb.BmMongodb
}

// GetByID according to the ID for the client information
func (p *PhAuthorizeCodeClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	in := PhModel.OpenClient{ID: id}
	out := PhModel.OpenClient{ID: id}
	err = p.mdb.FindOne(&in, &out)
	if err == nil {
		cli = Model2ClientInfo(&out)
		return
	}
	err = errors.New("not found by Id = " + id)
	return
}

// Set set client information
func (p *PhAuthorizeCodeClientStore) Set(cli oauth2.ClientInfo) (err error) {
	model := clientInfo2Model(cli).(*PhModel.OpenClient)
	_, err = p.mdb.InsertBmObject(model)
	return
}

func clientInfo2Model(cli oauth2.ClientInfo) (bmb BmModel.BmModelBase) {
	bmb = &PhModel.OpenClient{
		ID:     cli.GetID(),
		Secret: cli.GetSecret(),
		Domain: cli.GetDomain(),
		UserID: cli.GetUserID(),
	}
	return
}

func Model2ClientInfo(bmb BmModel.BmModelBase) (cli oauth2.ClientInfo) {
	model := bmb.(*PhModel.OpenClient)
	cli = &models.Client{
		ID:     model.ID,
		Secret: model.Secret,
		Domain: model.Domain,
		UserID: model.UserID,
	}
	return
}
