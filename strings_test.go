package parameters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCamelCaseToSnakeCase
func TestCamelCaseToSnakeCase(t *testing.T) {
	entries := map[string]string{
		"ID":          "id",
		"User":        "user",
		"UserName":    "user_name",
		"UserID":      "user_id",
		"MyJSON":      "my_json",
		"ProfileHTML": "profile_html",
		"RequestXML":  "request_xml",
	}

	for k, v := range entries {
		t.Run("testing: "+k, func(t *testing.T) {
			transformed := CamelToSnakeCase(k)
			assert.Equal(t, v, transformed)
		})
	}
}

// TestSnakeCaseToCamelCase
func TestSnakeCaseToCamelCase(t *testing.T) {
	entries := map[string]string{
		"id":           "ID",
		"user":         "User",
		"user_name":    "UserName",
		"user_id":      "UserID",
		"my_json":      "MyJSON",
		"profile_html": "ProfileHTML",
		"request_xml":  "RequestXML",
	}

	for k, v := range entries {
		t.Run("testing: "+k, func(t *testing.T) {
			transformed := SnakeToCamelCase(k, true)
			assert.Equal(t, v, transformed)
		})
	}
}
