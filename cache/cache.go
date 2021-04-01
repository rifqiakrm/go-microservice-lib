package cache

import (
	"fmt"
)


func Get(key string) ([]byte, error) {
	cacheKey := key

	bytes, err := GetCache(cacheKey)

	if err != nil {
		return nil, fmt.Errorf("error while retrieving cache : %v", err.Error())
	}

	return bytes, nil
}

func Set(key string, data interface{}, time int) error {
	cacheKey := key

	if err := SetCache(cacheKey, data, time); err != nil {
		return err
	}

	return nil
}
