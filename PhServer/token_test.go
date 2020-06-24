package PhServer

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v3/models"
	"github.com/PharbersDeveloper/PhAuthServer/PhUnits/yaml"
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

	convey.Convey("Test authorization code store", t, func() {
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
		convey.So(err, convey.ShouldBeNil)

		cinfo, err := tokenStore.GetByCode(info.Code)
		convey.So(err, convey.ShouldBeNil)
		convey.So(cinfo.GetUserID(), convey.ShouldEqual, info.UserID)

		err = tokenStore.RemoveByCode(info.Code)
		convey.So(err, convey.ShouldBeNil)

		cinfo, err = tokenStore.GetByCode(info.Code)
		convey.So(err, convey.ShouldBeNil)
		convey.So(cinfo, convey.ShouldBeNil)
	})

	convey.Convey("Test access token store", t, func() {
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
		convey.So(err, convey.ShouldBeNil)

		ainfo, err := tokenStore.GetByAccess(info.GetAccess())
		println(ainfo)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ainfo.GetUserID(), convey.ShouldEqual, info.GetUserID())

		err = tokenStore.RemoveByAccess(info.GetAccess())
		convey.So(err, convey.ShouldBeNil)

		ainfo, err = tokenStore.GetByAccess(info.GetAccess())
		println(ainfo)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ainfo, convey.ShouldBeNil)
	})

	convey.Convey("Test refresh token store", t, func() {
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
		convey.So(err, convey.ShouldBeNil)

		rinfo, err := tokenStore.GetByRefresh(info.GetRefresh())
		convey.So(err, convey.ShouldBeNil)
		convey.So(rinfo.GetUserID(), convey.ShouldEqual, info.GetUserID())

		err = tokenStore.RemoveByRefresh(info.GetRefresh())
		convey.So(err, convey.ShouldBeNil)

		rinfo, err = tokenStore.GetByRefresh(info.GetRefresh())
		convey.So(err, convey.ShouldBeNil)
		convey.So(rinfo, convey.ShouldBeNil)
	})

	convey.Convey("Test TTL", t, func() {
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
		convey.So(err, convey.ShouldBeNil)

		time.Sleep(time.Second * 1)
		ainfo, err := tokenStore.GetByAccess(info.Access)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ainfo, convey.ShouldBeNil)
		rinfo, err := tokenStore.GetByRefresh(info.Refresh)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rinfo, convey.ShouldBeNil)
	})
}
