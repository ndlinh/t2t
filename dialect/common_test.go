package dialect

import "testing"

func TestCamelCase(t *testing.T) {
	tests := map[string]string{
		"test_field": "TestField",
		"TestField":  "TestField",
		"testField":  "TestField",
		"api":        "API",
		"ip":         "IP",
		"Ip":         "IP",
		"IP":         "IP",
		"id":         "ID",
		"Id":         "ID",
		"ID":         "ID",
	}
	for k, v := range tests {
		if CamelCase(k) != v {
			t.Fatalf("%s failed", k)
		}
	}
}
