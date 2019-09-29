package parameters

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// TestGetParams_ParseJSONBody
func TestGetParams_ParseJSONBody(t *testing.T) {
	body := "{ \"test\": true }"
	r, err := http.NewRequest("POST", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}
}

// TestGetParams_ParseJSONBodyContentType
func TestGetParams_ParseJSONBodyContentType(t *testing.T) {
	body := "{ \"test\": true }"
	r, err := http.NewRequest("POST", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json; charset=utf8")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}
}

// TestGetParams_ParseNestedJSONBody
func TestGetParams_ParseNestedJSONBody(t *testing.T) {
	body := "{ \"test\": true, \"coord\": { \"lat\": 50.505, \"lon\": 10.101 }}"
	r, err := http.NewRequest("POST", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}

	val, present = params.Get("coord")
	if !present {
		t.Fatal("Key: 'coord' not found")
	}

	coord := val.(map[string]interface{})

	lat, present := coord["lat"]
	if !present {
		t.Fatal("Key: 'lat' not found")
	}
	if lat != 50.505 {
		t.Fatal("Value of 'lat' should be 50.505, got: ", lat)
	}

	lat, present = params.Get("coord.lat")
	if !present {
		t.Fatal("Nested Key: 'lat' not found")
	}
	if lat != 50.505 {
		t.Fatal("Value of 'lat' should be 50.505, got: ", lat)
	}

	lon, present := coord["lon"]
	if !present {
		t.Fatal("Key: 'lon' not found")
	}
	if lon != 10.101 {
		t.Fatal("Value of 'lon' should be 10.101, got: ", lon)
	}

	lon, present = params.Get("coord.lon")
	if !present {
		t.Fatal("Nested Key: 'lon' not found")
	}
	if lon != 10.101 {
		t.Fatal("Value of 'lon' should be 10.101, got: ", lon)
	}
}

// TestGetParams
func TestGetParams(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "test?test=true", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}
}

// TestParams_GetStringOk
func TestParams_GetStringOk(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=string", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetStringOk("test")
	if !ok {
		t.Fatal("failed getting string parameter", val, ok)
	} else if val != "string" {
		t.Fatal("failed getting string value", val, ok)
	}
}

// TestParams_GetBoolOk
func TestParams_GetBoolOk(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=true", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetBoolOk("test")
	if !ok {
		t.Fatal("failed getting bool parameter", val, ok)
	} else if !val {
		t.Fatal("failed getting bool value", val, ok)
	}
}

// TestParams_GetBytesOk
func TestParams_GetBytesOk(t *testing.T) {
	testBytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		testBytes[i] = 'a' + byte(i%26)
	}
	testString := string(testBytes)

	body := ""
	r, err := http.NewRequest("GET", "/test?test="+testString, strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetBytesOk("test")
	if !ok {
		t.Fatal("failed getting bytes parameter", val, ok)
	}
}

// TestParams_GetFloatOk
func TestParams_GetFloatOk(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123.1234", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetFloatOk("test")
	if !ok {
		t.Fatal("failed getting float parameter", val, ok)
	} else if val != 123.1234 {
		t.Fatal("failed getting float value", val, ok)
	}
}

// TestParams_GetIntOk
func TestParams_GetIntOk(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetIntOk("test")
	if !ok {
		t.Fatal("failed getting int parameter", val, ok)
	} else if val != 123 {
		t.Fatal("failed getting int value", val, ok)
	}
}

// TestParams_GetInt8Ok
func TestParams_GetInt8Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt8Ok("test")
	if !ok {
		t.Fatal("failed getting int parameter", val, ok)
	} else if val != 123 {
		t.Fatal("failed getting int value", val, ok)
	}
}

// TestParams_GetInt16Ok
func TestParams_GetInt16Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt16Ok("test")
	if !ok {
		t.Fatal("failed getting int parameter", val, ok)
	} else if val != 123 {
		t.Fatal("failed getting int value", val, ok)
	}
}

// TestParams_GetInt32Ok
func TestParams_GetInt32Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt32Ok("test")
	if !ok {
		t.Fatal("failed getting int parameter", val, ok)
	} else if val != 123 {
		t.Fatal("failed getting int value", val, ok)
	}
}

// TestParams_GetInt64Ok
func TestParams_GetInt64Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt64Ok("test")
	if !ok {
		t.Fatal("failed getting int parameter", val, ok)
	} else if val != 123 {
		t.Fatal("failed getting int value", val, ok)
	}
}

// TestParams_GetUint64Ok
func TestParams_GetUint64Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequest("GET", "/test?test=123", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetUint64Ok("test")
	if !ok {
		t.Fatal("failed getting uint parameter", val, ok)
	} else if val != 123 {
		t.Fatal("failed getting uint value", val, ok)
	}
}

// TestGetParams_Post
func TestGetParams_Post(t *testing.T) {
	body := "test=true"
	r, err := http.NewRequest("POST", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}
}

// TestGetParams_Put
func TestGetParams_Put(t *testing.T) {
	body := "test=true"
	r, err := http.NewRequest("PUT", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}
}

// TestGetParams_ParsePostUrlJSON
func TestGetParams_ParsePostUrlJSON(t *testing.T) {
	body := "{\"test\":true}"
	r, err := http.NewRequest("PUT", "test?test=false&id=1", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	if !present {
		t.Fatal("Key: 'test' not found")
	}
	if val != true {
		t.Fatal("Value of 'test' should be 'true', got: ", val)
	}

	val, present = params.GetFloatOk("id")
	if !present {
		t.Fatal("Key: 'id' not found")
	}
	if val != 1.0 {
		t.Fatal("Value of 'id' should be 1, got: ", val)
	}
}

// TestGetParams_ParseJSONBodyMux
func TestGetParams_ParseJSONBodyMux(t *testing.T) {
	body := "{ \"test\": true }"
	r, err := http.NewRequest("POST", "/test/42", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")
	m := mux.NewRouter()
	m.KeepContext = true
	m.HandleFunc("/test/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

		params := GetParams(r)

		val, present := params.Get("test")
		if !present {
			t.Fatal("Key: 'test' not found")
		}
		if val != true {
			t.Fatal("Value of 'test' should be 'true', got: ", val)
		}

		val, present = params.Get("id")
		if !present {
			t.Fatal("Key: 'id' not found")
		}
		if val != uint64(42) {
			t.Fatal("Value of 'id' should be 42, got: ", val)
		}
	})

	var match mux.RouteMatch
	if !m.Match(r, &match) {
		t.Error("Mux did not match")
	}
	m.ServeHTTP(nil, r)
}

// TestImbue
func TestImbue(t *testing.T) {
	body := "test=true&keys=this,that,something&values=1,2,3"
	r, err := http.NewRequest("PUT", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	type testType struct {
		Test   bool
		Keys   []string
		Values []int
	}

	var obj testType
	params.Imbue(&obj)

	if obj.Test != true {
		t.Fatal("Value of 'test' should be 'true', got: ", obj.Test)
	}
	if len(obj.Keys) != 3 {
		t.Fatal("Length of 'keys' should be '3', got: ", len(obj.Keys))
	}
	if len(obj.Values) != 3 {
		t.Fatal("Length of 'values' should be '3', got: ", len(obj.Values))
	}
	values := []int{1, 2, 3}
	for i, k := range obj.Values {
		if values[i] != k {
			t.Fatal("Expected ", values[i], ", got:", k)
		}
	}
}

// TestImbue_Time
func TestImbue_Time(t *testing.T) {
	body := "test=true&created_at=2016-06-07T00:30Z&remind_on=2016-07-17"
	r, err := http.NewRequest("PUT", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	type testType struct {
		Test      bool
		CreatedAt time.Time
		RemindOn  *time.Time
	}

	var obj testType
	params.Imbue(&obj)

	if obj.Test != true {
		t.Fatal("Value of 'test' should be 'true', got: ", obj.Test)
	}
	createdAt, _ := time.Parse(time.RFC3339, "2016-06-07T00:30Z00:00")
	if !obj.CreatedAt.Equal(createdAt) {
		t.Fatal("CreatedAt should be '2016-06-07T00:30Z', got:", obj.CreatedAt)
	}
	remindOn, _ := time.Parse(DateOnly, "2016-07-17")
	if obj.RemindOn == nil || !obj.RemindOn.Equal(remindOn) {
		t.Fatal("RemindOn should be '2016-07-17', got:", obj.RemindOn)
	}
}

// TestHasAll
func TestHasAll(t *testing.T) {
	body := "test=true&keys=this,that,something&values=1,2,3"
	r, err := http.NewRequest("PUT", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)
	//Test All
	if ok, missing := params.HasAll("test", "keys", "values"); !ok || len(missing) > 0 {
		t.Fatal("Params should have all keys, could not find", missing)
	}

	// Test Partial Contains
	if ok, missing := params.HasAll("test"); !ok || len(missing) > 0 {
		t.Fatal("Params should have key 'test', could not find", missing)
	}

	// Test Partial Missing
	if ok, missing := params.HasAll("test", "nope"); ok || len(missing) == 0 {
		t.Fatal("Params should not have key 'nope'", missing)
	} else if contains(missing, "test") {
		t.Fatal("Missing should not contain 'test'")
	}

	// Test All missing
	if ok, missing := params.HasAll("negative", "nope"); ok || len(missing) == 0 {
		t.Fatal("Params should not have key 'nope' nor 'negative'", missing)
	}
}

// TestGetParams_ParseEmpty test some garbage input, ids= "" (empty string) Should either be not ok, or empty slice
func TestGetParams_ParseEmpty(t *testing.T) {
	body := "{\"test\":true}"
	r, err := http.NewRequest("PUT", "test?ids=", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	ids, ok := params.GetUint64SliceOk("ids")
	if ok {
		if len(ids) > 0 {
			t.Fatal("ids should be !ok or an empty array. Length:", len(ids))
		}
	}
}

// TestGetParams_NegativeUint test Uint64 returns not ok for negative values
func TestGetParams_NegativeUint(t *testing.T) {
	body := "{\"id\":-1}"
	r, err := http.NewRequest("PUT", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	id, ok := params.GetUint64Ok("id")
	if ok {
		t.Fatal("Negative uint64 should be !ok not", id)
	}

	body = "{\"id\":1}"
	r, err = http.NewRequest("PUT", "test", strings.NewReader(body))
	if err != nil {
		t.Fatal("Could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params = GetParams(r)

	id, ok = params.GetUint64Ok("id")
	if !ok || id != 1 {
		t.Fatal("Id should be 1 not", id)
	}
}
