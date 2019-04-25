package PhServer

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v3/models"
	"ph_auth/PhUnits/yaml"
	"testing"
	"time"
)

func TestTokenStore(t *testing.T) {
	conf := yaml.LoadConfFromYAML("../resources/resource/service-def.yaml")
	args := conf.Daemons[1].Args
	db := BmRedis.BmRedis{
		Host:     args["host"],
		Port:     args["port"],
		Password: args["password"],
		Database: args["database"],
	}

	tokenStore, _ := NewAuthorizeCodeTokenStore(&db)

	Convey("Test authorization code store", t, func() {
		info := &models.Token{
			ClientID:      "1",
			UserID:        "1_1",
			RedirectURI:   "http://localhost/",
			Scope:         "all",
			Code:          "11_11_11",
			CodeCreateAt:  time.Now(),
			CodeExpiresIn: time.Second * 5,
		}
		err := tokenStore.Create(info)
		So(err, ShouldBeNil)

		cinfo, err := tokenStore.GetByCode(info.Code)
		So(err, ShouldBeNil)
		So(cinfo.GetUserID(), ShouldEqual, info.UserID)

		err = tokenStore.RemoveByCode(info.Code)
		So(err, ShouldBeNil)

		cinfo, err = tokenStore.GetByCode(info.Code)
		So(err, ShouldBeNil)
		So(cinfo, ShouldBeNil)
	})

	Convey("Test access token store", t, func() {
		info := &models.Token{
			ClientID:        "1",
			UserID:          "1_1",
			RedirectURI:     "http://localhost/",
			Scope:           "all",
			Access:          "1_1_1",
			AccessCreateAt:  time.Now(),
			AccessExpiresIn: time.Second * 5,
		}
		err := tokenStore.Create(info)
		So(err, ShouldBeNil)

		ainfo, err := tokenStore.GetByAccess(info.GetAccess())
		println(ainfo)
		So(err, ShouldBeNil)
		So(ainfo.GetUserID(), ShouldEqual, info.GetUserID())

		err = tokenStore.RemoveByAccess(info.GetAccess())
		So(err, ShouldBeNil)

		ainfo, err = tokenStore.GetByAccess(info.GetAccess())
		println(ainfo)
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)
	})

	Convey("Test refresh token store", t, func() {
		info := &models.Token{
			ClientID:         "1",
			UserID:           "1_2",
			RedirectURI:      "http://localhost/",
			Scope:            "all",
			Access:           "1_2_1",
			AccessCreateAt:   time.Now(),
			AccessExpiresIn:  time.Second * 5,
			Refresh:          "1_2_2",
			RefreshCreateAt:  time.Now(),
			RefreshExpiresIn: time.Second * 15,
		}
		err := tokenStore.Create(info)
		So(err, ShouldBeNil)

		rinfo, err := tokenStore.GetByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo.GetUserID(), ShouldEqual, info.GetUserID())

		err = tokenStore.RemoveByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)

		rinfo, err = tokenStore.GetByRefresh(info.GetRefresh())
		So(err, ShouldBeNil)
		So(rinfo, ShouldBeNil)
	})

	Convey("Test TTL", t, func() {
		info := &models.Token{
			ClientID:         "1",
			UserID:           "1_1",
			RedirectURI:      "http://localhost/",
			Scope:            "all",
			Access:           "1_3_1",
			AccessCreateAt:   time.Now(),
			AccessExpiresIn:  time.Second * 1,
			Refresh:          "1_3_2",
			RefreshCreateAt:  time.Now(),
			RefreshExpiresIn: time.Second * 1,
		}
		err := tokenStore.Create(info)
		So(err, ShouldBeNil)

		time.Sleep(time.Second * 1)
		ainfo, err := tokenStore.GetByAccess(info.Access)
		So(err, ShouldBeNil)
		So(ainfo, ShouldBeNil)
		rinfo, err := tokenStore.GetByRefresh(info.Refresh)
		So(err, ShouldBeNil)
		So(rinfo, ShouldBeNil)
	})
}
