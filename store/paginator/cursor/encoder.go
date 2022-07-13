// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cursor

import (
	"encoding/base64"
	"encoding/json"
	"megpoid.dev/go/contact-form/store/paginator/util"
	"reflect"
)

// NewEncoder creates cursor encoder
func NewEncoder(fields []EncoderField) *Encoder {
	return &Encoder{fields: fields}
}

// Encoder cursor encoder
type Encoder struct {
	fields []EncoderField
}

// EncoderField contains information about one encoder field.
type EncoderField struct {
	Key string
	// metas are needed for handling of custom types
	Meta any
}

// Encode encodes model into cursor
func (e *Encoder) Encode(model any) (string, error) {
	b, err := e.marshalJSON(model)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (e *Encoder) marshalJSON(model any) ([]byte, error) {
	rv := util.ReflectValue(model)
	fields := make([]any, len(e.fields))
	for i, field := range e.fields {
		f := rv.FieldByName(field.Key)
		if f == (reflect.Value{}) {
			return nil, ErrInvalidModel
		}
		if e.isNilable(f) && f.IsZero() {
			fields[i] = nil
		} else {
			// fetch values from custom types
			if ct, ok := util.ReflectValue(f).Interface().(CustomType); ok {
				value, err := ct.GetCustomTypeValue(field.Meta)
				if err != nil {
					return nil, err
				}
				fields[i] = value
			} else {
				fields[i] = util.ReflectValue(f).Interface()
			}
		}
	}
	result, err := json.Marshal(fields)
	if err != nil {
		return nil, ErrInvalidModel
	}
	return result, nil
}

func (e *Encoder) isNilable(v reflect.Value) bool {
	return v.Kind() >= 18 && v.Kind() <= 23
}
