package utils

import (
	"strconv"
	"strings"
)

type QueryFilter struct {
	store map[string]string
}

type FilterItem struct {
	value string
}

func NewQueryFilter(filter string) *QueryFilter {
	store := make(map[string]string)
	filterItems := strings.Split(filter, ",")
	for _, item := range filterItems {
		pair := strings.Split(item, ":")
		if len(pair) != 2 {
			continue
		}

		store[pair[0]] = pair[1]
	}

	return &QueryFilter{
		store: store,
	}
}

func (q *QueryFilter) Get(key string) (*FilterItem, bool) {
	if val, ok := q.store[key]; ok {
		return &FilterItem{value: val}, true
	}
	return nil, false
}

func (f *FilterItem) String() string {
	return f.value
}

func (f *FilterItem) Int64() int64 {
	i, _ := strconv.ParseInt(f.value, 10, 64)
	return i
}

func (f *FilterItem) Bool() bool {
	b, _ := strconv.ParseBool(f.value)
	return b
}
