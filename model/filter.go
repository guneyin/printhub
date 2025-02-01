package model

import (
	"fmt"
	"reflect"
	"strings"
)

type QueryFilter interface {
	Query() (interface{}, []interface{})
}

func Query(obj QueryFilter) (string, []interface{}) {
	var cond, arg string
	query := strings.Builder{}
	var args []interface{}

	vals := reflect.ValueOf(obj)
	for i := range vals.NumField() {
		field := vals.Field(i)

		if !field.IsZero() {
			param := vals.Type().Field(i).Tag.Get("query")
			if param == "" {
				continue
				//args = append(args, field.FieldByName(vals.Type().Field(i).Name).Interface())
			}

			params := strings.Split(param, ",")

			cond = "="
			arg = field.String()
			if len(params) > 1 {
				if params[1] == "like" {
					cond = "like"
					arg = fmt.Sprintf("%%%s%%", arg)
				}
			}

			query.WriteString(params[0])
			query.WriteString(cond)
			query.WriteString("?")

			args = append(args, arg)
		}
	}

	if query.Len() == 0 {
		query.WriteString("")
	}

	q := query.String()
	return q, args
}
