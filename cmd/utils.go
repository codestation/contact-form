// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

var unmarshalDecoder = viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	HexStringToByteArray(),
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
))

func HexStringToByteArray() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String || t != typeOfBytes {
			return data, nil
		}

		if strings.HasPrefix(data.(string), "s:") {
			return []byte(strings.Split(data.(string), ":")[1]), nil
		}

		return hex.DecodeString(data.(string))
	}
}
