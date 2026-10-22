package storage

import (
	"encoding/json"
	"hw2/internal/product"
	"os"
	"strconv"
)

type JsonStorage struct {
	filePath string
}

func (s *JsonStorage) LoadProducts() (map[int64]product.Product, error) {
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	unmarshalledDataMap := map[string]product.Product{}

	err = json.Unmarshal(fileData, &unmarshalledDataMap)
	if err != nil {
		return nil, err
	}

	productsMap := map[int64]product.Product{}
	for strKey, value := range unmarshalledDataMap {
		newIntKey, err := strconv.ParseInt(strKey, 10, 64)
		if err != nil {
			return nil, err
		}
		productsMap[newIntKey] = value
	}

	return productsMap, nil
}

func NewJsonStorage(filePath string) (*JsonStorage, error) {
	JsonStorage := JsonStorage{
		filePath: filePath,
	}
	return &JsonStorage, nil
}
