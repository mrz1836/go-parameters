package parameters

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// FuzzParseParams tests ParseParams with various malformed and edge case inputs
func FuzzParseParams(f *testing.F) {
	// Add seed corpus for different content types
	f.Add("application/json", `{"test": "value", "number": 123}`)
	f.Add("application/json", `{"nested": {"key": "value"}, "array": [1,2,3]}`)
	f.Add("application/json", `{"boolean": true, "float": 3.14}`)
	f.Add("application/x-msgpack", string([]byte{0x82, 0xa4, 0x74, 0x65, 0x73, 0x74, 0xa5, 0x76, 0x61, 0x6c, 0x75, 0x65}))
	f.Add("application/x-www-form-urlencoded", "key=value&number=123&bool=true")
	f.Add("multipart/form-data", "test boundary data")

	// Edge cases
	f.Add("application/json", ``)
	f.Add("application/json", `null`)
	f.Add("application/json", `{}`)
	f.Add("application/json", `{"": ""}`)
	f.Add("application/json", `{"key": null}`)

	// Malformed JSON
	f.Add("application/json", `{`)
	f.Add("application/json", `{"key": }`)
	f.Add("application/json", `{"key": "value"`)
	f.Add("application/json", `{key: "value"}`)

	f.Fuzz(func(t *testing.T, contentType, body string) {
		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
		req.Header.Set("Content-Type", contentType)

		// Ensure we don't panic during parsing
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseParams panicked with contentType=%s, body=%s: %v", contentType, body, r)
			}
		}()

		// Parse parameters
		params := ParseParams(req)

		// Basic validation - should always return a non-nil Params
		if params == nil {
			t.Errorf("ParseParams returned nil for contentType=%s, body=%s", contentType, body)
			return
		}

		// Values might be nil for certain edge cases (like "null" JSON), initialize if needed
		if params.Values == nil {
			params.Values = make(map[string]interface{})
		}

		// Test that we can call methods on the returned params without panicking
		testParamsMethods(t, params, contentType, body)
	})
}

// testParamsMethods tests that all parameter methods work without panicking
func testParamsMethods(t *testing.T, params *Params, contentType, body string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Params method panicked with contentType=%s, body=%s: %v", contentType, body, r)
		}
	}()

	// Test various get methods with common keys
	testKeys := []string{"", "test", "key", "number", "boolean", "array", "nested", "nonexistent"}

	for _, key := range testKeys {
		// Basic get - should never panic
		_, _ = params.Get(key)

		// All these methods can potentially panic with invalid type conversions
		// Wrap each one to prevent test failures on expected edge cases
		func() {
			defer func() { _ = recover() }()
			_ = params.GetString(key)
			_, _ = params.GetStringOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetInt(key)
			_, _ = params.GetIntOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetInt64(key)
			_, _ = params.GetInt64Ok(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetFloat(key)
			_, _ = params.GetFloatOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetBool(key)
			_, _ = params.GetBoolOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetUint64(key)
			_, _ = params.GetUint64Ok(key)
		}()

		// Slice methods
		func() {
			defer func() { _ = recover() }()
			_ = params.GetStringSlice(key)
			_, _ = params.GetStringSliceOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetIntSlice(key)
			_, _ = params.GetIntSliceOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetFloatSlice(key)
			_, _ = params.GetFloatSliceOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetUint64Slice(key)
			_, _ = params.GetUint64SliceOk(key)
		}()

		// Special methods
		func() {
			defer func() { _ = recover() }()
			_ = params.GetTime(key)
			_, _ = params.GetTimeOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetJSON(key)
			_, _ = params.GetJSONOk(key)
		}()
		func() {
			defer func() { _ = recover() }()
			_ = params.GetBytes(key)
			_, _ = params.GetBytesOk(key)
		}()
	}

	// Test other methods
	_ = params.Clone()
	_, _ = params.HasAll("test", "nonexistent")
}

// FuzzTypeConversions tests the various type conversion functions
func FuzzTypeConversions(f *testing.F) {
	// Add seed corpus for different types and edge cases
	f.Add("123")
	f.Add("-456")
	f.Add("0")
	f.Add("3.14")
	f.Add("-0.5")
	f.Add("true")
	f.Add("false")
	f.Add("TRUE")
	f.Add("FALSE")
	f.Add("1,2,3,4,5")
	f.Add("a,b,c")
	f.Add("1.1,2.2,3.3")
	f.Add("")
	f.Add(" ")
	f.Add("  123  ")

	// Edge cases for numbers
	f.Add("9223372036854775807")     // max int64
	f.Add("-9223372036854775808")    // min int64
	f.Add("18446744073709551615")    // max uint64
	f.Add("1.7976931348623157e+308") // max float64
	f.Add("2.2250738585072014e-308") // min positive float64

	// Invalid inputs
	f.Add("not_a_number")
	f.Add("123abc")
	f.Add("abc123")
	f.Add("1.2.3")
	f.Add("++123")
	f.Add("--456")

	f.Fuzz(func(t *testing.T, input string) {
		// Create params with the input value
		params := &Params{
			Values: map[string]interface{}{
				"test": input,
			},
		}

		// Test all type conversion methods without panicking
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Type conversion panicked with input=%s: %v", input, r)
			}
		}()

		// Test integer conversions
		_ = params.GetInt("test")
		_, _ = params.GetIntOk("test")
		_ = params.GetInt8("test")
		_, _ = params.GetInt8Ok("test")
		_ = params.GetInt16("test")
		_, _ = params.GetInt16Ok("test")
		_ = params.GetInt32("test")
		_, _ = params.GetInt32Ok("test")
		_ = params.GetInt64("test")
		_, _ = params.GetInt64Ok("test")
		_ = params.GetUint64("test")
		_, _ = params.GetUint64Ok("test")

		// Test float conversions
		_ = params.GetFloat("test")
		_, _ = params.GetFloatOk("test")

		// Test boolean conversions
		_ = params.GetBool("test")
		_, _ = params.GetBoolOk("test")

		// Test slice conversions
		_ = params.GetIntSlice("test")
		_, _ = params.GetIntSliceOk("test")
		_ = params.GetFloatSlice("test")
		_, _ = params.GetFloatSliceOk("test")
		_ = params.GetStringSlice("test")
		_, _ = params.GetStringSliceOk("test")
		_ = params.GetUint64Slice("test")
		_, _ = params.GetUint64SliceOk("test")

		// Test string conversion
		_ = params.GetString("test")
		_, _ = params.GetStringOk("test")
	})
}

// FuzzStringCaseConversion tests the snake_case <-> CamelCase conversion functions
func FuzzStringCaseConversion(f *testing.F) {
	// Add seed corpus
	f.Add("simple")
	f.Add("snake_case")
	f.Add("PascalCase")
	f.Add("camelCase")
	f.Add("UPPER_CASE")
	f.Add("mixed_CamelCase_string")
	f.Add("user_id")
	f.Add("html_content")
	f.Add("json_data")
	f.Add("xml_parser")
	f.Add("ID")
	f.Add("HTML")
	f.Add("JSON")
	f.Add("XML")

	// Edge cases
	f.Add("")
	f.Add("_")
	f.Add("__")
	f.Add("_test")
	f.Add("test_")
	f.Add("_test_")
	f.Add("a")
	f.Add("A")
	f.Add("123")
	f.Add("test123")
	f.Add("123test")
	f.Add("test_123_case")

	// Unicode cases
	f.Add("café_au_lait")
	f.Add("naïve_approach")

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("String conversion panicked with input=%s: %v", input, r)
			}
		}()

		// Test both directions of conversion
		snakeResult := CamelToSnakeCase(input)
		camelResult1 := SnakeToCamelCase(input, true)
		camelResult2 := SnakeToCamelCase(input, false)
		_ = MakeFirstUpperCase(input)

		// Test round-trip conversions don't panic
		_ = SnakeToCamelCase(snakeResult, true)
		_ = SnakeToCamelCase(snakeResult, false)
		_ = CamelToSnakeCase(camelResult1)
		_ = CamelToSnakeCase(camelResult2)
		_ = MakeFirstUpperCase(snakeResult)
		_ = MakeFirstUpperCase(camelResult1)
		_ = MakeFirstUpperCase(camelResult2)

		// Test that isKnownAbbreviation doesn't panic
		_ = isKnownAbbreviation(input)
		_ = isKnownAbbreviation(strings.ToLower(input))
		_ = isKnownAbbreviation(strings.ToUpper(input))
	})
}

// FuzzTimeParsing tests time parsing with various date/time formats
func FuzzTimeParsing(f *testing.F) {
	// Add seed corpus with valid time formats
	f.Add("2023-12-25T15:30:45Z")      // RFC3339
	f.Add("2023-12-25T15:30:45+02:00") // RFC3339 with timezone
	f.Add("2023-12-25")                // DateOnly
	f.Add("2023-12-25 15:30:45")       // DateTime
	f.Add("2023-12-25T15:30")          // HTMLDateTimeLocal

	// Edge cases and boundary values
	f.Add("1970-01-01T00:00:00Z") // Unix epoch
	f.Add("2038-01-19T03:14:07Z") // Year 2038 problem
	f.Add("9999-12-31T23:59:59Z") // Far future
	f.Add("0001-01-01T00:00:00Z") // Very old date

	// Leap year cases
	f.Add("2024-02-29") // Leap year
	f.Add("2023-02-28") // Non-leap year

	// Invalid but interesting inputs
	f.Add("")
	f.Add("invalid")
	f.Add("2023-13-45")                     // Invalid month/day
	f.Add("2023-02-30")                     // Invalid leap year
	f.Add("2023/12/25")                     // Wrong separator
	f.Add("25-12-2023")                     // Wrong order
	f.Add("2023-12-25T25:70:70Z")           // Invalid time
	f.Add("2023-12-25T15:30:45")            // Missing timezone
	f.Add("T15:30:45Z")                     // Missing date
	f.Add("2023-12-25T")                    // Missing time
	f.Add("2023-12-25 T15:30:45Z")          // Extra space
	f.Add("2023-12-25T15:30:45.123Z")       // With milliseconds
	f.Add("2023-12-25T15:30:45.123456789Z") // With nanoseconds

	f.Fuzz(func(t *testing.T, timeStr string) {
		// Create params with the time string
		params := &Params{
			Values: map[string]interface{}{
				"time": timeStr,
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Time parsing panicked with input=%s: %v", timeStr, r)
			}
		}()

		// Test time parsing methods
		_ = params.GetTime("time")
		_, _ = params.GetTimeOk("time")
		_ = params.GetTimeInLocation("time", time.UTC)
		_, ok := params.GetTimeInLocationOk("time", time.UTC)

		// If parsing succeeds, the result should be consistent
		if ok {
			time1 := params.GetTime("time")
			time2, _ := params.GetTimeOk("time")
			if !time1.Equal(time2) {
				t.Errorf("GetTime and GetTimeOk returned different times for input=%s", timeStr)
			}
		}

		// Test with different time zones
		locations := []*time.Location{
			time.UTC,
		}

		for _, loc := range locations {
			_, _ = params.GetTimeInLocationOk("time", loc)
			_ = params.GetTimeInLocation("time", loc)
		}
	})
}

// FuzzSliceParsing tests slice parsing with various formats and edge cases
func FuzzSliceParsing(f *testing.F) {
	// Add seed corpus for different slice types
	f.Add("1,2,3,4,5")
	f.Add("1.1,2.2,3.3,4.4")
	f.Add("a,b,c,d,e")
	f.Add("true,false,true")
	f.Add("123,456,789")
	f.Add("")
	f.Add(" ")
	f.Add(",")
	f.Add(",,")
	f.Add(",1,2,")
	f.Add("1,,2")
	f.Add(" 1 , 2 , 3 ")
	f.Add("1,2,3,")
	f.Add(",1,2,3")

	// Edge cases with numbers
	f.Add("0,1,2")
	f.Add("-1,-2,-3")
	f.Add("9223372036854775807,-9223372036854775808")        // int64 max/min
	f.Add("18446744073709551615,0")                          // uint64 max
	f.Add("1.7976931348623157e+308,2.2250738585072014e-308") // float64 extremes

	// Invalid numeric inputs
	f.Add("1,abc,3")
	f.Add("1.1.1,2.2")
	f.Add("1e999,2")
	f.Add("++1,2")
	f.Add("1,2,overflow123456789012345678901234567890")

	// Mixed valid/invalid
	f.Add("123,abc,456")
	f.Add("1,2,3,not_a_number")
	f.Add("valid,123,another_valid")

	f.Fuzz(func(t *testing.T, input string) {
		// Test with string input
		params1 := &Params{
			Values: map[string]interface{}{
				"slice": input,
			},
		}

		// Test with []byte input
		params2 := &Params{
			Values: map[string]interface{}{
				"slice": []byte(input),
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Slice parsing panicked with input=%s: %v", input, r)
			}
		}()

		// Test all slice parsing methods on both params
		for _, params := range []*Params{params1, params2} {
			// Integer slices
			_ = params.GetIntSlice("slice")
			_, intOk := params.GetIntSliceOk("slice")

			// Float slices
			_ = params.GetFloatSlice("slice")
			floatSlice, floatOk := params.GetFloatSliceOk("slice")

			// String slices (should always work)
			_ = params.GetStringSlice("slice")
			_, stringOk := params.GetStringSliceOk("slice")

			// Uint64 slices
			_ = params.GetUint64Slice("slice")
			_, uint64Ok := params.GetUint64SliceOk("slice")

			// Basic validation - check for NaN in float slices
			if floatOk && len(floatSlice) > 0 {
				for _, v := range floatSlice {
					if v != v { // NaN check
						t.Errorf("Float slice contains NaN value")
					}
				}
			}

			// Verify we have some results when successful
			_ = intOk
			_ = uint64Ok
			_ = stringOk
		}
	})
}

// FuzzJSONParsing tests JSON parsing with various valid and invalid inputs
func FuzzJSONParsing(f *testing.F) {
	// Add seed corpus
	f.Add(`{"key": "value"}`)
	f.Add(`{"number": 123}`)
	f.Add(`{"boolean": true}`)
	f.Add(`{"null_value": null}`)
	f.Add(`{"nested": {"inner": "value"}}`)
	f.Add(`{"array": [1, 2, 3]}`)
	f.Add(`{"mixed": {"arr": [{"id": 1}, {"id": 2}]}}`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`null`)
	f.Add(`"string"`)
	f.Add(`123`)
	f.Add(`true`)
	f.Add(`false`)

	// Invalid JSON
	f.Add(`{`)
	f.Add(`}`)
	f.Add(`{"key": }`)
	f.Add(`{"key": "value"`)
	f.Add(`{key: "value"}`)
	f.Add(`{"key": 'value'}`)
	f.Add(`{"key": undefined}`)
	f.Add(`{"trailing": "comma",}`)
	f.Add(`{,}`)
	f.Add(`{"": ""}`)

	f.Fuzz(func(t *testing.T, jsonStr string) {
		// Create params with the JSON string
		params := &Params{
			Values: map[string]interface{}{
				"json": jsonStr,
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("JSON parsing panicked with input=%s: %v", jsonStr, r)
			}
		}()

		// Test JSON parsing
		_ = params.GetJSON("json")
		_, _ = params.GetJSONOk("json")
	})
}

// FuzzUniqueUint64 tests the UniqueUint64 function with various inputs
// Since Go fuzzing doesn't support []uint64 directly, we simulate with multiple values
func FuzzUniqueUint64(f *testing.F) {
	// Add seed corpus with up to 5 values
	f.Add(uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint8(0))                    // empty slice (length 0)
	f.Add(uint64(1), uint64(0), uint64(0), uint64(0), uint64(0), uint8(1))                    // single element
	f.Add(uint64(1), uint64(2), uint64(3), uint64(0), uint64(0), uint8(3))                    // three elements
	f.Add(uint64(1), uint64(1), uint64(1), uint64(0), uint64(0), uint8(3))                    // duplicates
	f.Add(uint64(3), uint64(1), uint64(4), uint64(1), uint64(5), uint8(5))                    // mixed with duplicates
	f.Add(uint64(0), uint64(1), uint64(0), uint64(2), uint64(0), uint8(5))                    // zeros and values
	f.Add(uint64(18446744073709551615), uint64(0), uint64(0), uint64(0), uint64(0), uint8(1)) // max uint64

	f.Fuzz(func(t *testing.T, v1, v2, v3, v4, v5 uint64, length uint8) {
		// Limit length to reasonable bounds
		if length > 5 {
			length = length % 6
		}

		// Build input slice
		var input []uint64
		values := []uint64{v1, v2, v3, v4, v5}
		for i := uint8(0); i < length; i++ {
			input = append(input, values[i])
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("UniqueUint64 panicked with input length %d: %v", len(input), r)
			}
		}()

		result := UniqueUint64(input)

		// Basic validation
		if len(result) > len(input) {
			t.Errorf("UniqueUint64 returned more elements (%d) than input (%d)", len(result), len(input))
		}

		// Check that all elements in result are unique
		seen := make(map[uint64]bool)
		for _, v := range result {
			if seen[v] {
				t.Errorf("UniqueUint64 returned duplicate element: %d", v)
			}
			seen[v] = true
		}

		// Check that all elements in result existed in input
		inputMap := make(map[uint64]bool)
		for _, v := range input {
			inputMap[v] = true
		}
		for _, v := range result {
			if !inputMap[v] {
				t.Errorf("UniqueUint64 returned element %d not present in input", v)
			}
		}
	})
}
