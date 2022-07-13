// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"megpoid.dev/go/contact-form/store/paginator/util"
	"reflect"
)

// Rule for paginator
type Rule struct {
	Key             string
	Order           Order
	SQLRepr         string
	CustomType      *CustomType
	NULLReplacement any
}

// CustomType for paginator. It provides extra info needed to paginate across custom types (e.g. JSON)
type CustomType struct {
	Meta interface{}
	Type reflect.Type
}

func (r *Rule) validate(dest any) (err error) {
	if _, ok := util.ReflectType(dest).FieldByName(r.Key); !ok {
		return ErrInvalidModel
	}
	if r.Order != "" {
		if err = r.Order.validate(); err != nil {
			return
		}
	}
	return nil
}
