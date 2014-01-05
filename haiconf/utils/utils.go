// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

func IsDir(d string) (bool, error) {
	fi, err := os.Stat(d)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func HasFileName(f string) bool {
	ext := filepath.Ext(f)
	return len(ext) > 0
}

func ToStringMap(m map[string]interface{}) (map[string]string, error) {
	value := reflect.ValueOf(m)
	if value.Kind() != reflect.Map {
		return map[string]string{}, errors.New("Not a map received")
	}

	keys := value.MapKeys()
	stringMap := make(map[string]string, len(keys))
	for _, key := range keys {
		stringMap[key.String()] = value.MapIndex(key).Elem().String()
	}

	return stringMap, nil
}
