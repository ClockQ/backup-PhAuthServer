package PhServer

import (
	"time"
	"encoding/json"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/PharbersDeveloper/PhAuthServer/PhUnits/uuid"
)

// NewAuthorizeCodeTokenStore create a token store instance based on redis
func NewAuthorizeCodeTokenStore(rdb *BmRedis.BmRedis) (store oauth2.TokenStore, err error) {
	store = &PhAuthorizeCodeTokenStore{
		rdb: rdb,
	}
	return
}

// PhAuthorizeCodeTokenStore token storage based on buntdb(https://github.com/tidwall/buntdb)
// PhAuthorizeCodeTokenStore implement ==> oauth2.TokenStore
type PhAuthorizeCodeTokenStore struct {
	rdb *BmRedis.BmRedis
}

// Create create and store the new token information
func (ts *PhAuthorizeCodeTokenStore) Create(info oauth2.TokenInfo) (err error) {
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}

	client := ts.rdb.GetRedisClient()
	defer client.Close()

	if code := info.GetCode(); code != "" {
		client.Set(code, string(jv), info.GetCodeExpiresIn())
		return
	} else {
		basicID := uuid.Must(uuid.NewRandom()).String()
		aexp := info.GetAccessExpiresIn()
		rexp := aexp
		if refresh := info.GetRefresh(); refresh != "" {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
			client.Set(refresh, basicID, rexp)
		}

		client.Set(basicID, string(jv), rexp)
		client.Set(info.GetAccess(), basicID, aexp)

		return
	}
}

// remove key
func (ts *PhAuthorizeCodeTokenStore) remove(key string) (err error) {
	//verr := ts.db.Update(func(tx *buntdb.Tx) (err error) {
	//	_, err = tx.Delete(key)
	//	return
	//})
	//if verr == buntdb.ErrNotFound {
	//	return
	//}
	//err = verr
	return
}

// RemoveByCode use the authorization code to delete the token information
func (ts *PhAuthorizeCodeTokenStore) RemoveByCode(code string) (err error) {
	//err = ts.remove(code)
	return
}

// RemoveByAccess use the access token to delete the token information
func (ts *PhAuthorizeCodeTokenStore) RemoveByAccess(access string) (err error) {
	//err = ts.remove(access)
	return
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *PhAuthorizeCodeTokenStore) RemoveByRefresh(refresh string) (err error) {
	//err = ts.remove(refresh)
	return
}

func (ts *PhAuthorizeCodeTokenStore) getData(key string) (ti oauth2.TokenInfo, err error) {
	client := ts.rdb.GetRedisClient()
	defer client.Close()

	jv, err := client.Get(key).Result()
	if err != nil {
		return
	}
	var tm models.Token
	err = json.Unmarshal([]byte(jv), &tm)
	if err != nil {
		return
	}
	ti = &tm
	return
}

func (ts *PhAuthorizeCodeTokenStore) getBasicID(key string) (basicID string, err error) {
	//verr := ts.db.View(func(tx *buntdb.Tx) (err error) {
	//	v, err := tx.Get(key)
	//	if err != nil {
	//		return
	//	}
	//	basicID = v
	//	return
	//})
	//if verr != nil {
	//	if verr == buntdb.ErrNotFound {
	//		return
	//	}
	//	err = verr
	//}
	return
}

// GetByCode use the authorization code for token information data
func (ts *PhAuthorizeCodeTokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.getData(code)
	return
}

// GetByAccess use the access token for token information data
func (ts *PhAuthorizeCodeTokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	//basicID, err := ts.getBasicID(access)
	//if err != nil {
	//	return
	//}
	//ti, err = ts.getData(basicID)
	return
}

// GetByRefresh use the refresh token for token information data
func (ts *PhAuthorizeCodeTokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	//basicID, err := ts.getBasicID(refresh)
	//if err != nil {
	//	return
	//}
	//ti, err = ts.getData(basicID)
	return
}
