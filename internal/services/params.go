package services

import (
	"fmt"
	"net/url"
)

type Param func(values url.Values)

func DefaultParam(name string, value string) Param {
	return func(values url.Values) {
		values.Add(name, value)
	}
}

func Limit(limit int) Param {
	return DefaultParam("limit", fmt.Sprint(limit))
}

func Offset(offset int) Param {
	return DefaultParam("offset", fmt.Sprint(offset))
}

func SortBy(value string) Param {
	return DefaultParam("sort_by", value)
}

func OrderBy(value string) Param {
	return DefaultParam("order_by", value)
}

func ApplyParams(values url.Values, params []Param) {
	for _, p := range params {
		p(values)
	}
}
