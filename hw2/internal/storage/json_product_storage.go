package storage

import (
	"encoding/json"
	"fmt"
	"hw2/internal/product"
	"os"
	"strconv"
)

type JsonProductStorage struct {
	filePath string
}

func (s *JsonProductStorage) LoadProducts() (map[int64]product.Product, error) {
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", s.filePath, err)
	}

	unmarshalledDataMap := map[string]product.Product{}

	err = json.Unmarshal(fileData, &unmarshalledDataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	productsMap := map[int64]product.Product{}
	for strKey, value := range unmarshalledDataMap {
		newIntKey, err := strconv.ParseInt(strKey, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sku %q: %w", strKey, err)
		}
		productsMap[newIntKey] = value
	}

	return productsMap, nil
}

func NewJsonProductStorage(filePath string) (*JsonProductStorage, error) {
	jsonStorage := JsonProductStorage{
		filePath: filePath,
	}
	return &jsonStorage, nil
}
