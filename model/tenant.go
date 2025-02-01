package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model `json:"-"`
	UUID       string `json:"uuid" gorm:"index"`
	Email      string `json:"email" gorm:"uniqueIndex"`
	Name       string `json:"name" gorm:"index"`
	Address    string `json:"address"`
	Logo       string `json:"logo"`
}

type TenantList []Tenant

type TenantFilter struct {
	ID   string `query:"uuid"`
	Name string `query:"name,like"`
}

func (tf TenantFilter) Query() (interface{}, []interface{}) {
	return Query(tf)
	//query := interface{}(nil)
	//var cond, arg string
	//query := strings.Builder{}
	//var args []interface{}
	//
	//vals := reflect.ValueOf(&tf).Elem()
	//for i := range vals.NumField() {
	//	field := vals.Field(i)
	//
	//	if !field.IsZero() {
	//		param := vals.Type().Field(i).Tag.Get("query")
	//		if param == "" {
	//			continue
	//			//args = append(args, field.FieldByName(vals.Type().Field(i).Name).Interface())
	//		}
	//
	//		params := strings.Split(param, ",")
	//		//query.WriteString(params[0])
	//		//query.WriteString("=")
	//		//query.WriteString("?")
	//
	//		cond = "="
	//		arg = field.String()
	//		if len(params) > 1 {
	//			if params[1] == "like" {
	//				cond = "like"
	//				arg = fmt.Sprintf("%%%s%%", arg)
	//			}
	//		}
	//
	//		query.WriteString(params[0])
	//		query.WriteString(cond)
	//		query.WriteString("?")
	//
	//		args = append(args, arg)
	//	}
	//}
	//
	//return query, args
}

type TenantUser struct {
	gorm.Model
	TenantID uint `gorm:"index"`
	UserID   uint
}

//type TenantUserList []UserList

func NewTenant(d []byte) (*Tenant, error) {
	t := &Tenant{}
	err := json.Unmarshal(d, t)
	return t, err
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	t.UUID = uuid.New().String()

	return nil
}
