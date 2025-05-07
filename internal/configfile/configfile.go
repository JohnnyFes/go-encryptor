package configfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var debugMode bool

// SetDebug включает или выключает debug-режим
func SetDebug(enabled bool) {
	debugMode = enabled
}

func debugPrint(format string, args ...interface{}) {
	if debugMode {
		fmt.Printf(format, args...)
	}
}

// UpdateConfigFile обновляет указанные поля в YAML/JSON файле зашифрованными значениями
func UpdateConfigFile(configPath string, fields, values []string) error {
	if len(fields) != len(values) {
		return fmt.Errorf("number of fields and values must match")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	var root map[string]interface{}
	ext := filepath.Ext(configPath)
	isYAML := ext == ".yaml" || ext == ".yml"
	if isYAML {
		var raw map[interface{}]interface{}
		err = yaml.Unmarshal(data, &raw)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
		root = convertMapI2MapS(raw)
	} else {
		err = json.Unmarshal(data, &root)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
	}
	// Если root nil (например, пустой файл), инициализируем map
	if root == nil {
		root = make(map[string]interface{})
	}
	// Обновляем поля
	for i, f := range fields {
		if err := setNestedField(root, f, values[i]); err != nil {
			return fmt.Errorf("failed to set field %s: %w", f, err)
		}
	}
	// ВРЕМЕННО: выводим root map и абсолютный путь к файлу для отладки
	absPath, _ := filepath.Abs(configPath)
	debugPrint("[DEBUG] Writing to: %s\n", absPath)
	debugPrint("[DEBUG] Data to write: %#v\n", root)
	// Сохраняем обратно
	var out []byte
	if isYAML {
		out, err = yaml.Marshal(root)
	} else {
		out, err = json.MarshalIndent(root, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// setNestedField устанавливает значение вложенного поля по пути вида "a.b.c"
func setNestedField(m map[string]interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	last := len(parts) - 1
	cur := m
	for i, p := range parts {
		if i == last {
			cur[p] = value
			return nil
		}
		if next, ok := cur[p].(map[string]interface{}); ok {
			cur = next
		} else if cur[p] == nil {
			next := make(map[string]interface{})
			cur[p] = next
			cur = next
		} else {
			return fmt.Errorf("field %s is not a map", p)
		}
	}
	return nil
}

// convertMapI2MapS рекурсивно преобразует map[interface{}]interface{} в map[string]interface{}
func convertMapI2MapS(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(m))
	for k, v := range m {
		ks, ok := k.(string)
		if !ok {
			continue
		}
		switch vv := v.(type) {
		case map[interface{}]interface{}:
			res[ks] = convertMapI2MapS(vv)
		case []interface{}:
			res[ks] = convertSliceI2SliceS(vv)
		default:
			res[ks] = vv
		}
	}
	return res
}

// convertSliceI2SliceS рекурсивно преобразует элементы с map[interface{}]interface{} внутри слайса
func convertSliceI2SliceS(s []interface{}) []interface{} {
	for i, v := range s {
		switch vv := v.(type) {
		case map[interface{}]interface{}:
			s[i] = convertMapI2MapS(vv)
		case []interface{}:
			s[i] = convertSliceI2SliceS(vv)
		}
	}
	return s
}
