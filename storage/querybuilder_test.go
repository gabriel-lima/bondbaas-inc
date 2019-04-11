package storage

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateInsertQuery(t *testing.T) {
	var tests = []struct {
		tableName string
		payload   map[string]interface{}
		expected  string
	}{
		{
			"products", map[string]interface{}{"name": "Coca-cola"},
			`INSERT INTO products (id, name) VALUES (DEFAULT, $1)`,
		},
		{
			"customers", map[string]interface{}{"name": "Coca-cola", "price": 3.5},
			`INSERT INTO customers (id, name, price) VALUES (DEFAULT, $1, $2)`,
		},
		{
			"users", map[string]interface{}{},
			`INSERT INTO users (id) VALUES (DEFAULT)`,
		},
		{
			"groups", map[string]interface{}{"id": 123},
			`INSERT INTO groups (id) VALUES (DEFAULT)`,
		},
	}

	for _, test := range tests {
		if output := GenerateInsertQuery(test.tableName, test.payload); output != test.expected {
			t.Errorf("Test Failed:\nExpected: %s\nReceived: %s", test.expected, output)
		}
	}
}

func TestGenerateInsertValues(t *testing.T) {
	var testCases = []struct {
		payload  map[string]interface{}
		expected []interface{}
	}{
		{
			map[string]interface{}{"name": "Coca-cola"},
			[]interface{}{"Coca-cola"},
		},
		{
			map[string]interface{}{"name": "Pepsi", "price": 3.5},
			[]interface{}{"Pepsi", 3.5},
		},
		{
			map[string]interface{}{"id": 12345},
			[]interface{}{12345},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.payload), func(t *testing.T) {
			if output := GenerateInsertValues(tc.payload); !reflect.DeepEqual(output, tc.expected) {
				t.Errorf("Test Failed:\nExpected: %v\nReceived: %v", tc.expected, output)
			}
		})
	}
}

func TestGenerateUpdateQuery(t *testing.T) {
	var tests = []struct {
		tableName string
		payload   map[string]interface{}
		expected  string
	}{
		{
			"products", map[string]interface{}{"name": "Coca-cola"},
			`UPDATE products SET name = $2 WHERE id = $1`,
		},
		{
			"customers", map[string]interface{}{"name": "Coca-cola", "price": 3.5},
			`UPDATE customers SET name = $2, price = $3 WHERE id = $1`,
		},
		{
			"users", map[string]interface{}{"id": 12345},
			``,
		},
		{
			"groups", map[string]interface{}{"id": 12345, "name": "Pepsi"},
			`UPDATE groups SET name = $2 WHERE id = $1`,
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
