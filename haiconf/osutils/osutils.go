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
