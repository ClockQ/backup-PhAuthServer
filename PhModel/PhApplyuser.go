package PhModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Applyuser struct {
	ID         string        `json:"-"`
	Id_        bson.ObjectId `json:"-" bson:"_id"`
	Name       string        `json:"name" bson:"name"`
	Company    string        `json:"company" bson:"company"`
	Email      string        `json:"email" bson:"email"`
	Phone      string        `json:"phone" bson:"phone"`
	Department string        `json:"department" bson:"department"` // 新增字段 部门
	Position   string        `json:"position" bson:"position"`     // 新增字段 职务
	Intention  string        `json:"intention" bson:"intention"`   // 新增字段 感兴趣内容
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Applyuser) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Applyuser) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Applyuser) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	r := make(map[string]interface{})
	var ids []bson.ObjectId
	for k, v := range parameters {
		switch k {
		case "ids":
			for i := 0; i < len(v); i++ {
				ids = append(ids, bson.ObjectIdHex(v[i]))
			}
			r["$in"] = ids
			rst["_id"] = r
		case "email":
			rst[k] = v[0]
		}
	}
	return rst
}