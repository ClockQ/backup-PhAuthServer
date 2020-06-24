package PhServer

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v3/models"
	"github.com/PharbersDeveloper/PhAuthServer/PhUnits/yaml"
	"testing"
)

func TestClientStore(t *testing.T) {
	conf := yaml.LoadConfFromYAML("../resources/resource/service-def.yaml")
	args := conf.Daemons[0].Args
	db := BmMongodb.BmMongodb{
		Host:     args["host"],
		Port:     args["port"],
		Database: args["database"],
	}

	clientStore := NewAuthorizeCodeClientStore(&db)

	Convey("Test client store", t, func() {
		id, err := clientStore.SetInfo(&models.Client{Secret: "test"})
		So(err, ShouldBeNil)

		cli, err := clientStore.GetByID(id)
		So(err, ShouldBeNil)
		So(cli.GetID(), ShouldEqual, id)

		err = clientStore.DelInfo(&models.Client{ID: id})
		So(err, ShouldBeNil)
	})
}
