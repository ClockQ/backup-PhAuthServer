package PhServer

import (
	"errors"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmModel"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"ph_auth/PhModel"
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
	if !bson.IsObjectIdHex(id) {
		err = errors.New(id + " isn't ObjectIdHex")
		return
	}

	in := PhModel.Client{ID: id}
	out := PhModel.Client{ID: id}
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
	_, err = p.SetInfo(cli)
	return
}

// SetInfo SetInfo call by Set
func (p *PhAuthorizeCodeClientStore) SetInfo(cli oauth2.ClientInfo) (id string, err error) {
	model := clientInfo2Model(cli).(*PhModel.Client)
	id, err = p.mdb.InsertBmObject(model)
	return
}

// DelInfo
func (p *PhAuthorizeCodeClientStore) DelInfo(cli oauth2.ClientInfo) (err error) {
	model := clientInfo2Model(cli).(*PhModel.Client)
	err = p.mdb.Delete(model)
	return
}

func clientInfo2Model(cli oauth2.ClientInfo) (bmb BmModel.BmModelBase) {
	bmb = &PhModel.Client{
		ClientID:  cli.GetID(),
		Secret:    cli.GetSecret(),
		Domain:    cli.GetDomain(),
		AccountID: cli.GetUserID(),
	}
	return
}

func Model2ClientInfo(bmb BmModel.BmModelBase) (cli oauth2.ClientInfo) {
	model := bmb.(*PhModel.Client)
	cli = &models.Client{
		ID:     model.ClientID,
		Secret: model.Secret,
		Domain: model.Domain,
		UserID: model.AccountID,
	}
	return
}
