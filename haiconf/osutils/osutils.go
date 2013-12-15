// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package osutils

import (
	"os"
)

func GetEnvList(enVars []string) map[string]string {
	tmpBucket := make(map[string]string, len(enVars))

	for _, name := range enVars {
		v := os.Getenv(name)
		if v != "" {
			tmpBucket[name] = v
		}
	}

	return tmpBucket
}

func SetEnvList(enVars map[string]string) error {
	var err error

	for name, value := range enVars {
		err = os.Setenv(name, value)
		if err != nil {
			return err
		}
	}

	return nil
}
