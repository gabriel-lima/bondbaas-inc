package storage

import (
	"reflect"
	"testing"
)

func TestGenerateUpdateQuery(t *testing.T) {
	var tests = []struct {
		tableName string
		payload   map[string]interface{}
		expected  string
	}{
		{
			"product", map[string]interface{}{"name": "Coca-cola"},
			`UPDATE product SET name = $2 WHERE id = $1`,
		},
		{
			"product", map[string]interface{}{"name": "Coca-cola", "price": 3.5},
			`UPDATE product SET name = $2, price = $3 WHERE id = $1`,
		},
		{
			"product", map[string]interface{}{"id": 12345},
			``,
		},
		{
			"product", map[string]interface{}{"id": 12345, "name": "Pepsi"},
			`UPDATE product SET name = $2 WHERE id = $1`,
		},
	}

	for _, test := range tests {
		if output := GenerateUpdateQuery(test.tableName, test.payload); output != test.expected {
			t.Errorf("Test Failed:\nExpected: %s\nReceived: %s", test.expected, output)
		}
	}
}

func TestGenerateUpdateValues(t *testing.T) {
	var tests = []struct {
		ID       int
		payload  map[string]interface{}
		expected []interface{}
	}{
		{
			123,
			map[string]interface{}{"name": "Coca-cola"},
			[]interface{}{123, "Coca-cola"},
		},
		{
			1234,
			map[string]interface{}{"name": "Pepsi", "price": 3.5},
			[]interface{}{1234, "Pepsi", 3.5},
		},
		{
			12345,
			map[string]interface{}{"id": 12345},
			[]interface{}{12345, 12345},
		},
	}

	for _, test := range tests {
		if output := GenerateUpdateValues(test.ID, test.payload); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test Failed:\nExpected: %v\nReceived: %v", test.expected, output)
		}
	}
}
