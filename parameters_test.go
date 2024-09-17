package parameters

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const testJSONParam = `{ "test": true }`

// TestGetParams_ParseJSONBody tests the method with JSON body
func TestGetParams_ParseJSONBody(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "test", strings.NewReader(testJSONParam))
	assert.NoError(t, err)
	assert.NotNil(t, r)

	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// BenchmarkGetParams_ParseJSONBody benchmarks the method
func BenchmarkGetParams_ParseJSONBody(b *testing.B) {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "test", strings.NewReader(testJSONParam))
	assert.NoError(b, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	for i := 0; i < b.N; i++ {
		_ = GetParams(r)
	}
}

// TestGetParams_ParseJSONBodyContentType tests the method with JSON body and content type
func TestGetParams_ParseJSONBodyContentType(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "test", strings.NewReader(testJSONParam))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json; charset=utf8")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// TestGetParams_ParseNestedJSONBody tests the method with nested JSON
func TestGetParams_ParseNestedJSONBody(t *testing.T) {
	body := "{ \"test\": true, \"coordinate\": { \"lat\": 50.505, \"lon\": 10.101 }}"
	r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)

	val, present = params.Get("coordinate")
	assert.Equal(t, true, present)

	coordinate := val.(map[string]interface{})

	var lat interface{}
	lat, present = coordinate["lat"]
	assert.Equal(t, true, present)
	assert.Equal(t, 50.505, lat)

	lat, present = params.Get("coordinate.lat")
	assert.Equal(t, true, present)
	assert.Equal(t, 50.505, lat)

	var lon interface{}
	lon, present = coordinate["lon"]
	assert.Equal(t, true, present)
	assert.Equal(t, 10.101, lon)

	lon, present = params.Get("coordinate.lon")
	assert.Equal(t, true, present)
	assert.Equal(t, 10.101, lon)
}

// TestGetParams tests the GetParams method
func TestGetParams(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "test?test=true", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// BenchmarkGetParams benchmarks the method
func BenchmarkGetParams(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "test?test=true", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	for i := 0; i < b.N; i++ {
		_ = GetParams(r)
	}
}

// TestParams_GetStringOk tests the GetStringOk method
func TestParams_GetStringOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=string", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetStringOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, "string", val)

	val = params.GetString("test")
	assert.Equal(t, "string", val)
}

// BenchmarkGetStringOk benchmarks the method
func BenchmarkParams_GetStringOk(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=string", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetStringOk("test")
	}
}

// TestParams_GetBoolOk tests the GetBoolOk method
func TestParams_GetBoolOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=true", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetBoolOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, true, val)
}

// BenchmarkGetBoolOk benchmarks the method
func BenchmarkParams_GetBoolOk(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=true", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetBoolOk("test")
	}
}

// TestParams_GetBytesOk tests the GetBytesOk method
func TestParams_GetBytesOk(t *testing.T) {
	testBytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		testBytes[i] = 'a' + byte(i%26)
	}
	testString := string(testBytes)

	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test="+testString, strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetBytesOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, 75, len(val))

	val = params.GetBytes("test")
	assert.Equal(t, 75, len(val))
}

// BenchmarkParams_GetBytesOk benchmarks the method
func BenchmarkParams_GetBytesOk(b *testing.B) {
	testBytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		testBytes[i] = 'a' + byte(i%26)
	}
	testString := string(testBytes)

	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test="+testString, strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetBytesOk("test")
	}
}

// BenchmarkParams_GetBool benchmarks the method
func BenchmarkParams_GetBool(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=true", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_ = params.GetBool("test")
	}
}

// TestParams_GetFloatOk tests the GetFloatOk method
func TestParams_GetFloatOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123.1234", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetFloatOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, 123.1234, val)

	val = params.GetFloat("test")
	assert.Equal(t, 123.1234, val)
}

// BenchmarkParams_GetFloatOk benchmarks the method
func BenchmarkParams_GetFloatOk(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123.1234", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetFloatOk("test")
	}
}

// TestParams_GetFloatOk_Zero tests the GetFloatOk method
func TestParams_GetFloatOk_Zero(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=null", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetFloatOk("test")
	assert.Equal(t, float64(0), val)
	assert.Equal(t, true, ok)
}

// TestParams_GetIntOk tests the GetIntOk method
func TestParams_GetIntOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetIntOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, 123, val)

	val = params.GetInt("test")
	assert.Equal(t, 123, val)
}

// BenchmarkParams_GetIntOk benchmarks the method
func BenchmarkParams_GetIntOk(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetIntOk("test")
	}
}

// TestParams_GetInt8Ok tests the GetInt8 method
func TestParams_GetInt8Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt8Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int8(123), val)
}

// TestParams_GetInt8TooSmall tests the GetInt8 method
func TestParams_GetInt8TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=-300", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt8("test")
	assert.Equal(t, int8(0), val)
}

// TestParams_GetInt8TooBig tests the GetInt8 method
func TestParams_GetInt8TooBig(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=300", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt8("test")
	assert.Equal(t, int8(0), val)
}

// TestParams_GetInt16Ok tests the GetInt16 method
func TestParams_GetInt16Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt16Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int16(123), val)

	val = params.GetInt16("test")
	assert.Equal(t, int16(123), val)
}

// TestParams_GetInt16TooSmall tests the GetInt16 method
func TestParams_GetInt16TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=-32769", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt16("test")
	assert.Equal(t, int16(0), val)
}

// TestParams_GetInt16TooBig tests the GetInt16 method
func TestParams_GetInt16TooBig(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=32769", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt16("test")
	assert.Equal(t, int16(0), val)
}

// TestParams_GetInt32Ok tests the GetInt32 method
func TestParams_GetInt32Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt32Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int32(123), val)

	val = params.GetInt32("test")
	assert.Equal(t, int32(123), val)
}

// TestParams_GetInt32TooSmall tests the GetInt32 method
func TestParams_GetInt32TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("/test?test=%d", math.MinInt32-1), strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt32("test")
	assert.Equal(t, int32(0), val)
}

// TestParams_GetInt32TooBig tests the GetInt32 method
func TestParams_GetInt32TooBig(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("/test?test=%d", math.MaxInt32+1), strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt32("test")
	assert.Equal(t, int32(0), val)
}

// TestParams_GetInt64Ok tests the GetInt64 method
func TestParams_GetInt64Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt64Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(123), val)

	val = params.GetInt64("test")
	assert.Equal(t, int64(123), val)
}

// BenchmarkParams_GetInt64Ok benchmarks the method
func BenchmarkParams_GetInt64Ok(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetInt64Ok("test")
	}
}

// TestParams_GetInt64TooSmall tests the GetInt64 method
func TestParams_GetInt64TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("/test?test=%d", 0), strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt64("test")
	assert.Equal(t, int64(0), val)
}

// TestParams_GetUint64Ok tests the GetUint64Ok method
func TestParams_GetUint64Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetUint64Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, uint64(123), val)
}

// TestGetParams_Post tests the method with a POST request
func TestGetParams_Post(t *testing.T) {
	body := "test=true"
	r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// TestParams_GetTimeOk tests the GetTimeOk method
func TestParams_GetTimeOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=2020-12-31", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetTimeOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, "2020-12-31 00:00:00 +0000 UTC", val.String())
}

// TestGetParams_Put tests the method with a PUT request
func TestGetParams_Put(t *testing.T) {
	body := "test=true"
	r, err := http.NewRequestWithContext(context.Background(), http.MethodPut, "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// TestGetParams_ParsePostUrlJSON tests the method with a JSON body
func TestGetParams_ParsePostUrlJSON(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test?test=false&id=1", strings.NewReader(testJSONParam))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)

	val, present = params.GetFloatOk("id")
	assert.Equal(t, true, present)
	assert.Equal(t, 1.0, val)
}

// TestGetParams_ParseJSONBodyMux tests the method with mux
func TestGetParams_ParseJSONBodyMux(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/test/42", strings.NewReader(testJSONParam))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")
	m := mux.NewRouter()
	// m.KeepContext = true
	m.HandleFunc("/test/{id:[0-9]+}", func(_ http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

		params := GetParams(r)

		val, present := params.Get("test")
		assert.Equal(t, true, present)
		assert.Equal(t, true, val)

		val, present = params.Get("id")
		assert.Equal(t, true, present)
		assert.Equal(t, uint64(42), val)
	})

	var match mux.RouteMatch
	assert.Equal(t, true, m.Match(r, &match))
	m.ServeHTTP(nil, r)
}

// TestImbue tests the Imbue method
func TestImbue(t *testing.T) {
	body := "test=true&keys=this,that,something&values=1,2,3"
	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test", strings.NewReader(body))
	assert.NoError(t, err)
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

	assert.Equal(t, true, obj.Test)
	assert.Equal(t, 3, len(obj.Keys))
	assert.Equal(t, 3, len(obj.Values))

	values := []int{1, 2, 3}
	for i, k := range obj.Values {
		assert.Equal(t, k, values[i])
	}
}

// TestImbue_Time tests the Imbue method with time.Time
func TestImbue_Time(t *testing.T) {
	body := "test=true&created_at=2016-06-07T00:30Z&remind_on=2016-07-17"
	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test", strings.NewReader(body))
	assert.NoError(t, err)
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

	assert.Equal(t, true, obj.Test)

	createdAt, _ := time.Parse(time.RFC3339, "2016-06-07T00:30Z00:00")
	assert.Equal(t, true, obj.CreatedAt.Equal(createdAt))

	remindOn, _ := time.Parse(DateOnly, "2016-07-17")
	assert.NotNil(t, remindOn)
	assert.Equal(t, true, obj.RemindOn.Equal(remindOn))
}

// TestHasAll tests the HasAll method
func TestHasAll(t *testing.T) {
	body := "test=true&keys=this,that,something&values=1,2,3"
	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	t.Run("test all", func(t *testing.T) {
		ok, missing := params.HasAll("test", "keys", "values")
		assert.Equal(t, true, ok)
		assert.Equal(t, 0, len(missing))
	})

	t.Run("test partial", func(t *testing.T) {
		ok, missing := params.HasAll("test")
		assert.Equal(t, true, ok)
		assert.Equal(t, 0, len(missing))
	})

	t.Run("test partial missing", func(t *testing.T) {
		ok, missing := params.HasAll("test", "nope")
		assert.Equal(t, false, ok)
		assert.NotEqual(t, 0, len(missing))
	})

	t.Run("test all missing", func(t *testing.T) {
		ok, missing := params.HasAll("negative", "nope")
		assert.Equal(t, false, ok)
		assert.NotEqual(t, 0, len(missing))
	})
}

// TestGetParams_ParseEmpty test some garbage input, ids= "" (empty string) Should either be not ok, or empty slice
func TestGetParams_ParseEmpty(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test?ids=", strings.NewReader(testJSONParam))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	ids, ok := params.GetUint64SliceOk("ids")
	assert.Equal(t, true, ok)
	assert.Equal(t, 0, len(ids))
}

// TestGetParams_NegativeUint test Uint64 returns not ok for negative values
func TestGetParams_NegativeUint(t *testing.T) {
	body := "{\"id\":-1}"
	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	id, ok := params.GetUint64Ok("id")
	assert.Equal(t, false, ok)
	assert.Equal(t, uint64(0), id)

	body = "{\"id\":1}"
	r, err = http.NewRequestWithContext(context.Background(), "PUT", "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params = GetParams(r)

	id, ok = params.GetUint64Ok("id")
	assert.Equal(t, true, ok)
	assert.Equal(t, uint64(1), id)
}

// TestNestedStructs (from go-parameters:master)
func TestNestedStructs(t *testing.T) {
	type testStruct struct {
		Val        uint64 `json:"val"`
		NestStruct struct {
			Field1 string `json:"field_1"`
		} `json:"nest_struct"`
	}

	body := `{
		"val": 1234,
		"nest_struct": {
			"field_1": "Hello World"
		}
	}`

	expected := &testStruct{
		Val: 1234,
		NestStruct: struct {
			Field1 string `json:"field_1"`
		}{
			Field1: "Hello World",
		},
	}

	r, err := http.NewRequestWithContext(
		context.Background(), http.MethodPut, "test", strings.NewReader(body),
	)
	if err != nil {
		t.Fatal("could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	testObj := &testStruct{}
	params.Imbue(testObj)

	if !reflect.DeepEqual(testObj, expected) {
		t.Fatalf("expected %+v, Got %+v", expected, testObj)
	}
}

// TestCustomTypeSetter (from go-parameters:master)
func TestCustomTypeSetter(t *testing.T) {
	type testStruct struct {
		Val        uint64 `json:"val"`
		NestStruct struct {
			Field1 string `json:"field_1"`
		} `json:"nest_struct"`
	}

	body := `{
		"val": 1234,
		"nest_struct": {
			"field_1": "Hello World"
		}
	}`

	expected := &testStruct{
		Val: 1234,
		NestStruct: struct {
			Field1 string `json:"field_1"`
		}{
			Field1: "Goodbye World",
		},
	}
	CustomTypeSetter = func(field *reflect.Value, _ interface{}) error {
		if field.Type() == reflect.TypeOf(expected.NestStruct) {
			field.Set(reflect.ValueOf(expected.NestStruct))
			return nil
		}
		return errors.New("no type definition found")
	}

	r, err := http.NewRequestWithContext(
		context.Background(), http.MethodPut, "test", strings.NewReader(body),
	)
	if err != nil {
		t.Fatal("could not build request", err)
	}
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	testObj := &testStruct{}
	params.Imbue(testObj)

	if !reflect.DeepEqual(testObj, expected) {
		t.Fatalf("expected %+v, Got %+v", expected, testObj)
	}
}

// TestMakeParsedReq will test the MakeParsedReq function
func TestMakeParsedReq(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "test?test=true", strings.NewReader(""))
	assert.NoError(t, err)
	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))
	params := GetParams(r)

	fn := func(_ http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, params))
		p := GetParams(r)
		val, present := p.Get("test")
		assert.Equal(t, true, present)
		assert.Equal(t, true, val)
	}

	req := MakeParsedReq(fn)
	assert.NotNil(t, req)

	assert.NotPanics(t, func() {
		req(nil, r)
	})
}

// TestParams_GetFloatSliceOk tests the GetFloatSliceOk method
func TestParams_GetFloatSliceOk(t *testing.T) {
	tests := []struct {
		name           string
		params         *Params
		key            string
		expectedSlice  []float64
		expectedResult bool
	}{
		{
			name: "Value is []float64",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []float64{1.1, 2.2, 3.3},
				},
			},
			key:            "floats",
			expectedSlice:  []float64{1.1, 2.2, 3.3},
			expectedResult: true,
		},
		{
			name: "Value is comma-separated string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": "4.4,5.5,6.6",
				},
			},
			key:            "floats",
			expectedSlice:  []float64{4.4, 5.5, 6.6},
			expectedResult: true,
		},
		{
			name: "Value is []interface{} with float64 and string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []interface{}{7.7, "8.8", 9.9},
				},
			},
			key:            "floats",
			expectedSlice:  []float64{7.7, 8.8, 9.9},
			expectedResult: true,
		},
		{
			name: "Value is []interface{} with invalid string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []interface{}{10.1, "invalid", 11.1},
				},
			},
			key:            "floats",
			expectedSlice:  []float64{},
			expectedResult: false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:            "missing_key",
			expectedSlice:  []float64{},
			expectedResult: false,
		},
		{
			name: "Value is of unexpected type (int)",
			params: &Params{
				Values: map[string]interface{}{
					"floats": 123,
				},
			},
			key:            "floats",
			expectedSlice:  []float64{},
			expectedResult: false,
		},
		{
			name: "Value is empty string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": "",
				},
			},
			key:            "floats",
			expectedSlice:  []float64{},
			expectedResult: true,
		},
		{
			name: "Value is empty []interface{}",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []interface{}{},
				},
			},
			key:            "floats",
			expectedSlice:  []float64{},
			expectedResult: true,
		},
		{
			name: "Value is empty []float64",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []float64{},
				},
			},
			key:            "floats",
			expectedSlice:  []float64{},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice, result := tt.params.GetFloatSliceOk(tt.key)
			assert.Equal(t, tt.expectedResult, result, "Result mismatch")
			assert.Equal(t, tt.expectedSlice, slice, "Slice mismatch")
		})
	}
}

// TestParams_GetInt64Ok_Extended tests the GetInt64Ok method
func TestParams_GetInt64Ok_Extended(t *testing.T) {
	tests := []struct {
		name           string
		params         *Params
		key            string
		expectedValue  int64
		expectedResult bool
	}{
		{
			name: "Value is int within int64 range",
			params: &Params{
				Values: map[string]interface{}{
					"intKey": 12345,
				},
			},
			key:            "intKey",
			expectedValue:  12345,
			expectedResult: true,
		},
		{
			name: "Value is string representation of int",
			params: &Params{
				Values: map[string]interface{}{
					"strIntKey": "67890",
				},
			},
			key:            "strIntKey",
			expectedValue:  67890,
			expectedResult: true,
		},
		{
			name: "Value is float64 that can be converted to int64",
			params: &Params{
				Values: map[string]interface{}{
					"floatKey": 123.0,
				},
			},
			key:            "floatKey",
			expectedValue:  123,
			expectedResult: true,
		},
		{
			name: "Value is string that represents a float",
			params: &Params{
				Values: map[string]interface{}{
					"strFloatKey": "456.789",
				},
			},
			key:            "strFloatKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is of unexpected type (bool)",
			params: &Params{
				Values: map[string]interface{}{
					"boolKey": true,
				},
			},
			key:            "boolKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:            "missingKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value causes overflow",
			params: &Params{
				Values: map[string]interface{}{
					"overflowKey": "9223372036854775808", // math.MaxInt64 + 1
				},
			},
			key:            "overflowKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is negative",
			params: &Params{
				Values: map[string]interface{}{
					"negativeKey": -98765,
				},
			},
			key:            "negativeKey",
			expectedValue:  -98765,
			expectedResult: true,
		},
		{
			name: "Value is zero",
			params: &Params{
				Values: map[string]interface{}{
					"zeroKey": 0,
				},
			},
			key:            "zeroKey",
			expectedValue:  0,
			expectedResult: true,
		},
		{
			name: "Value is []byte representation of int",
			params: &Params{
				Values: map[string]interface{}{
					"byteKey": []byte("54321"),
				},
			},
			key:            "byteKey",
			expectedValue:  54321,
			expectedResult: true,
		},
		{
			name: "Value is invalid string",
			params: &Params{
				Values: map[string]interface{}{
					"invalidStrKey": "not_a_number",
				},
			},
			key:            "invalidStrKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is float64 with decimal part",
			params: &Params{
				Values: map[string]interface{}{
					"floatDecimalKey": 123.456,
				},
			},
			key:            "floatDecimalKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is max int64",
			params: &Params{
				Values: map[string]interface{}{
					"maxInt64Key": strconv.FormatInt(math.MaxInt64, 10),
				},
			},
			key:            "maxInt64Key",
			expectedValue:  math.MaxInt64,
			expectedResult: true,
		},
		{
			name: "Value is min int64",
			params: &Params{
				Values: map[string]interface{}{
					"minInt64Key": strconv.FormatInt(math.MinInt64, 10),
				},
			},
			key:            "minInt64Key",
			expectedValue:  math.MinInt64,
			expectedResult: true,
		},
		{
			name: "Value causes negative overflow",
			params: &Params{
				Values: map[string]interface{}{
					"negOverflowKey": "-9223372036854775809", // math.MinInt64 - 1
				},
			},
			key:            "negOverflowKey",
			expectedValue:  0,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, result := tt.params.GetInt64Ok(tt.key)
			assert.Equal(t, tt.expectedResult, result, "Result mismatch")
			assert.Equal(t, tt.expectedValue, value, "Value mismatch")
		})
	}
}

// TestParams_GetIntSliceOk tests the GetIntSliceOk method
func TestParams_GetIntSliceOk(t *testing.T) {
	tests := []struct {
		name           string
		params         *Params
		key            string
		expectedSlice  []int
		expectedResult bool
	}{
		{
			name: "Value is []int",
			params: &Params{
				Values: map[string]interface{}{
					"integers": []int{1, 2, 3},
				},
			},
			key:            "integers",
			expectedSlice:  []int{1, 2, 3},
			expectedResult: true,
		},
		{
			name: "Value is comma-separated string",
			params: &Params{
				Values: map[string]interface{}{
					"integers": "4,5,6",
				},
			},
			key:            "integers",
			expectedSlice:  []int{4, 5, 6},
			expectedResult: true,
		},
		{
			name: "Value is []byte of comma-separated integers",
			params: &Params{
				Values: map[string]interface{}{
					"integers": []byte("7,8,9"),
				},
			},
			key:            "integers",
			expectedSlice:  []int{7, 8, 9},
			expectedResult: true,
		},
		{
			name: "Value is []interface{} with integers",
			params: &Params{
				Values: map[string]interface{}{
					"integers": []interface{}{10, 11, 12},
				},
			},
			key:            "integers",
			expectedSlice:  []int{10, 11, 12},
			expectedResult: true,
		},
		{
			name: "Value is []interface{} with strings",
			params: &Params{
				Values: map[string]interface{}{
					"integers": []interface{}{"13", "14", "15"},
				},
			},
			key:            "integers",
			expectedSlice:  []int{13, 14, 15},
			expectedResult: true,
		},
		{
			name: "Value is []interface{} with float64",
			params: &Params{
				Values: map[string]interface{}{
					"integers": []interface{}{16.0, 17.0, 18.0},
				},
			},
			key:            "integers",
			expectedSlice:  []int{16, 17, 18},
			expectedResult: true,
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"integers": nil,
				},
			},
			key:            "integers",
			expectedSlice:  nil,
			expectedResult: true,
		},
		{
			name: "Value is empty string",
			params: &Params{
				Values: map[string]interface{}{
					"integers": "",
				},
			},
			key:            "integers",
			expectedSlice:  nil,
			expectedResult: true,
		},
		{
			name: "Value is of unexpected type (bool)",
			params: &Params{
				Values: map[string]interface{}{
					"integers": true,
				},
			},
			key:            "integers",
			expectedSlice:  nil,
			expectedResult: true,
		},
		{
			name: "Value contains invalid data in string",
			params: &Params{
				Values: map[string]interface{}{
					"integers": "19,invalid,20",
				},
			},
			key:            "integers",
			expectedSlice:  []int{19},
			expectedResult: false,
		},
		{
			name: "Value contains invalid data in []interface{}",
			params: &Params{
				Values: map[string]interface{}{
					"integers": []interface{}{21, "invalid", 22},
				},
			},
			key:            "integers",
			expectedSlice:  []int{21},
			expectedResult: false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:            "missing_key",
			expectedSlice:  []int{},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice, result := tt.params.GetIntSliceOk(tt.key)
			assert.Equal(t, tt.expectedResult, result, "Result mismatch")
			assert.Equal(t, tt.expectedSlice, slice, "Slice mismatch")
		})
	}
}

// BenchmarkParams_GetIntSliceOk benchmarks the GetIntSliceOk method
func BenchmarkParams_GetIntSliceOk(b *testing.B) {
	body := "integers=1,2,3,4,5"
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "test", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetIntSliceOk("integers")
	}
}

// TestParams_GetUint64Ok_Extended tests the GetUint64Ok method
func TestParams_GetUint64Ok_Extended(t *testing.T) {
	tests := []struct {
		name           string
		params         *Params
		key            string
		expectedValue  uint64
		expectedResult bool
	}{
		{
			name: "Value is uint64",
			params: &Params{
				Values: map[string]interface{}{
					"uint64Key": uint64(123456789),
				},
			},
			key:            "uint64Key",
			expectedValue:  123456789,
			expectedResult: true,
		},
		{
			name: "Value is string representing non-negative integer",
			params: &Params{
				Values: map[string]interface{}{
					"strUintKey": "987654321",
				},
			},
			key:            "strUintKey",
			expectedValue:  987654321,
			expectedResult: true,
		},
		{
			name: "Value is string representing negative integer",
			params: &Params{
				Values: map[string]interface{}{
					"strNegIntKey": "-12345",
				},
			},
			key:            "strNegIntKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is int64 (positive)",
			params: &Params{
				Values: map[string]interface{}{
					"int64PosKey": int64(12345),
				},
			},
			key:            "int64PosKey",
			expectedValue:  12345,
			expectedResult: true,
		},
		{
			name: "Value is int64 (negative)",
			params: &Params{
				Values: map[string]interface{}{
					"int64NegKey": int64(-67890),
				},
			},
			key:            "int64NegKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is float64 (positive)",
			params: &Params{
				Values: map[string]interface{}{
					"float64PosKey": 123.0,
				},
			},
			key:            "float64PosKey",
			expectedValue:  123,
			expectedResult: true,
		},
		{
			name: "Value is float64 (negative)",
			params: &Params{
				Values: map[string]interface{}{
					"float64NegKey": -456.0,
				},
			},
			key:            "float64NegKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is []byte representing non-negative integer",
			params: &Params{
				Values: map[string]interface{}{
					"byteKey": []byte("654321"),
				},
			},
			key:            "byteKey",
			expectedValue:  654321,
			expectedResult: true,
		},
		{
			name: "Value is of unexpected type (bool)",
			params: &Params{
				Values: map[string]interface{}{
					"boolKey": true,
				},
			},
			key:            "boolKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:            "missingKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value causes overflow",
			params: &Params{
				Values: map[string]interface{}{
					"overflowKey": strconv.FormatUint(math.MaxUint64, 10),
				},
			},
			key:            "overflowKey",
			expectedValue:  math.MaxUint64,
			expectedResult: true,
		},
		{
			name: "Value is zero",
			params: &Params{
				Values: map[string]interface{}{
					"zeroKey": 0,
				},
			},
			key:            "zeroKey",
			expectedValue:  0,
			expectedResult: true,
		},
		{
			name: "Value is string representing float",
			params: &Params{
				Values: map[string]interface{}{
					"strFloatKey": "123.456",
				},
			},
			key:            "strFloatKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is string representing negative float",
			params: &Params{
				Values: map[string]interface{}{
					"strNegFloatKey": "-789.123",
				},
			},
			key:            "strNegFloatKey",
			expectedValue:  0,
			expectedResult: false,
		},
		{
			name: "Value is uint (positive)",
			params: &Params{
				Values: map[string]interface{}{
					"uintKey": uint(55555),
				},
			},
			key:            "uintKey",
			expectedValue:  55555,
			expectedResult: true,
		},
		{
			name: "Value is string representing number larger than uint64",
			params: &Params{
				Values: map[string]interface{}{
					"overflowKey": "18446744073709551616", // math.MaxUint64 + 1
				},
			},
			key:            "overflowKey",
			expectedValue:  0,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, result := tt.params.GetUint64Ok(tt.key)
			assert.Equal(t, tt.expectedResult, result, "Result mismatch")
			assert.Equal(t, tt.expectedValue, value, "Value mismatch")
		})
	}
}

// BenchmarkParams_GetUint64Ok benchmarks the GetUint64Ok method
func BenchmarkParams_GetUint64Ok(b *testing.B) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test?test=123", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetUint64Ok("test")
	}
}

// TestParams_GetIntSliceOk_WithOverflow tests the GetIntSliceOk method with overflow
func TestParams_GetIntSliceOk_WithOverflow(t *testing.T) {
	var overflowValue string
	if strconv.IntSize == 32 {
		overflowValue = "2147483648" // math.MaxInt32 + 1
	} else {
		overflowValue = "9223372036854775808" // math.MaxInt64 + 1
	}

	params := &Params{
		Values: map[string]interface{}{
			"integers": overflowValue,
		},
	}

	slice, ok := params.GetIntSliceOk("integers")

	assert.False(t, ok, "Expected ok to be false due to overflow")
	assert.Equal(t, []int{}, slice, "Slice mismatch")
}

// TestParams_GetIntOk_WithOverflow tests the GetIntOk method with overflow
func TestParams_GetIntOk_WithOverflow(t *testing.T) {
	var overflowValue string
	if strconv.IntSize == 32 {
		overflowValue = "2147483648" // math.MaxInt32 + 1
	} else {
		overflowValue = "9223372036854775808" // math.MaxInt64 + 1
	}

	params := &Params{
		Values: map[string]interface{}{
			"intKey": overflowValue,
		},
	}

	expectedValue := 0
	value, ok := params.GetIntOk("intKey")

	assert.False(t, ok, "Expected ok to be false due to overflow")
	assert.Equal(t, expectedValue, value, "Value mismatch")
}

// TestParams_GetFloatSlice tests the GetFloatSlice method
func TestParams_GetFloatSlice(t *testing.T) {
	tests := []struct {
		name          string
		params        *Params
		key           string
		expectedSlice []float64
	}{
		{
			name: "Value is []float64",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []float64{1.1, 2.2, 3.3},
				},
			},
			key:           "floats",
			expectedSlice: []float64{1.1, 2.2, 3.3},
		},
		{
			name: "Value is comma-separated string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": "4.4,5.5,6.6",
				},
			},
			key:           "floats",
			expectedSlice: []float64{4.4, 5.5, 6.6},
		},
		{
			name: "Value is []byte of comma-separated numbers",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []byte("7.7,8.8,9.9"),
				},
			},
			key:           "floats",
			expectedSlice: []float64{7.7, 8.8, 9.9},
		},
		{
			name: "Value is []interface{} with float64s",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []interface{}{10.1, 11.2, 12.3},
				},
			},
			key:           "floats",
			expectedSlice: []float64{10.1, 11.2, 12.3},
		},
		{
			name: "Value is []interface{} with numeric strings",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []interface{}{"13.4", "14.5", "15.6"},
				},
			},
			key:           "floats",
			expectedSlice: []float64{13.4, 14.5, 15.6},
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"floats": nil,
				},
			},
			key:           "floats",
			expectedSlice: []float64{},
		},
		{
			name: "Value is empty string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": "",
				},
			},
			key:           "floats",
			expectedSlice: []float64{},
		},
		{
			name: "Value is of unexpected type (bool)",
			params: &Params{
				Values: map[string]interface{}{
					"floats": true,
				},
			},
			key:           "floats",
			expectedSlice: []float64{},
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:           "missing_key",
			expectedSlice: []float64{},
		},
		{
			name: "Value contains invalid data in string",
			params: &Params{
				Values: map[string]interface{}{
					"floats": "16.7,invalid,18.9",
				},
			},
			key:           "floats",
			expectedSlice: []float64{},
		},
		{
			name: "Value contains invalid data in []interface{}",
			params: &Params{
				Values: map[string]interface{}{
					"floats": []interface{}{19.0, "invalid", 20.1},
				},
			},
			key:           "floats",
			expectedSlice: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := tt.params.GetFloatSlice(tt.key)
			assert.Equal(t, tt.expectedSlice, slice, "Slice mismatch")
		})
	}
}

// TestParams_GetIntOk_Extended tests the GetIntOk method
func TestParams_GetIntOk_Extended(t *testing.T) {
	tests := []struct {
		name          string
		params        *Params
		key           string
		expectedValue int
		expectedOk    bool
	}{
		{
			name: "Value is int (positive)",
			params: &Params{
				Values: map[string]interface{}{
					"intKey": 123,
				},
			},
			key:           "intKey",
			expectedValue: 123,
			expectedOk:    true,
		},
		{
			name: "Value is int (negative)",
			params: &Params{
				Values: map[string]interface{}{
					"intKey": -456,
				},
			},
			key:           "intKey",
			expectedValue: -456,
			expectedOk:    true,
		},
		{
			name: "Value is int (zero)",
			params: &Params{
				Values: map[string]interface{}{
					"intKey": 0,
				},
			},
			key:           "intKey",
			expectedValue: 0,
			expectedOk:    true,
		},
		{
			name: "Value is int64 (positive)",
			params: &Params{
				Values: map[string]interface{}{
					"int64Key": int64(789),
				},
			},
			key:           "int64Key",
			expectedValue: 789,
			expectedOk:    true,
		},
		{
			name: "Value is int64 (negative)",
			params: &Params{
				Values: map[string]interface{}{
					"int64Key": int64(-1011),
				},
			},
			key:           "int64Key",
			expectedValue: -1011,
			expectedOk:    true,
		},
		{
			name: "Value is uint (within int range)",
			params: &Params{
				Values: map[string]interface{}{
					"uintKey": uint(2022),
				},
			},
			key:           "uintKey",
			expectedValue: 2022,
			expectedOk:    true,
		},
		{
			name: "Value is uint64 (exceeds int range)",
			params: &Params{
				Values: map[string]interface{}{
					"uint64Key": uint64(math.MaxInt) + 1,
				},
			},
			key:           "uint64Key",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is float64 (integer value)",
			params: &Params{
				Values: map[string]interface{}{
					"floatKey": 3033.0,
				},
			},
			key:           "floatKey",
			expectedValue: 3033,
			expectedOk:    true,
		},
		{
			name: "Value is float64 (non-integer value)",
			params: &Params{
				Values: map[string]interface{}{
					"floatKey": 4044.5,
				},
			},
			key:           "floatKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is string representing integer",
			params: &Params{
				Values: map[string]interface{}{
					"strKey": "5055",
				},
			},
			key:           "strKey",
			expectedValue: 5055,
			expectedOk:    true,
		},
		{
			name: "Value is string representing negative integer",
			params: &Params{
				Values: map[string]interface{}{
					"strKey": "-6066",
				},
			},
			key:           "strKey",
			expectedValue: -6066,
			expectedOk:    true,
		},
		{
			name: "Value is string representing number exceeding int range",
			params: &Params{
				Values: map[string]interface{}{
					"strKey": func() string {
						if strconv.IntSize == 32 {
							return "2147483648" // math.MaxInt32 + 1
						}
						return "9223372036854775808" // math.MaxInt64 + 1
					}(),
				},
			},
			key:           "strKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is string representing non-integer number",
			params: &Params{
				Values: map[string]interface{}{
					"strKey": "7077.8",
				},
			},
			key:           "strKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is []byte representing integer",
			params: &Params{
				Values: map[string]interface{}{
					"byteKey": []byte("8088"),
				},
			},
			key:           "byteKey",
			expectedValue: 8088,
			expectedOk:    true,
		},
		{
			name: "Value is of unexpected type (bool)",
			params: &Params{
				Values: map[string]interface{}{
					"boolKey": true,
				},
			},
			key:           "boolKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:           "missingKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"nilKey": nil,
				},
			},
			key:           "nilKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is math.MaxInt",
			params: &Params{
				Values: map[string]interface{}{
					"maxIntKey": math.MaxInt,
				},
			},
			key:           "maxIntKey",
			expectedValue: math.MaxInt,
			expectedOk:    true,
		},
		{
			name: "Value is math.MinInt",
			params: &Params{
				Values: map[string]interface{}{
					"minIntKey": math.MinInt,
				},
			},
			key:           "minIntKey",
			expectedValue: math.MinInt,
			expectedOk:    true,
		},
		{
			name: "Value is float64 exceeding int range",
			params: &Params{
				Values: map[string]interface{}{
					"floatKey": float64(math.MaxInt) + 1.0,
				},
			},
			key:           "floatKey",
			expectedValue: 0,
			expectedOk:    false,
		},
		{
			name: "Value is negative float64 integer",
			params: &Params{
				Values: map[string]interface{}{
					"floatKey": -9099.0,
				},
			},
			key:           "floatKey",
			expectedValue: -9099,
			expectedOk:    true,
		},
		{
			name: "Value is negative float64 non-integer",
			params: &Params{
				Values: map[string]interface{}{
					"floatKey": -10010.5,
				},
			},
			key:           "floatKey",
			expectedValue: 0,
			expectedOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := tt.params.GetIntOk(tt.key)
			assert.Equal(t, tt.expectedOk, ok, "Expected ok to be %v, got %v", tt.expectedOk, ok)
			assert.Equal(t, tt.expectedValue, value, "Expected value to be %v, got %v", tt.expectedValue, value)
		})
	}
}

// TestParams_GetStringSliceOk tests the GetStringSliceOk method
func TestParams_GetStringSliceOk(t *testing.T) {
	tests := []struct {
		name          string
		params        *Params
		key           string
		expectedSlice []string
		expectedOk    bool
	}{
		{
			name: "Value is []string",
			params: &Params{
				Values: map[string]interface{}{
					"key": []string{"a", "b", "c"},
				},
			},
			key:           "key",
			expectedSlice: []string{"a", "b", "c"},
			expectedOk:    true,
		},
		{
			name: "Value is string (comma-separated)",
			params: &Params{
				Values: map[string]interface{}{
					"key": "a,b,c",
				},
			},
			key:           "key",
			expectedSlice: []string{"a", "b", "c"},
			expectedOk:    true,
		},
		{
			name: "Value is []byte (comma-separated)",
			params: &Params{
				Values: map[string]interface{}{
					"key": []byte("a,b,c"),
				},
			},
			key:           "key",
			expectedSlice: []string{"a", "b", "c"},
			expectedOk:    true,
		},
		{
			name: "Value is []interface{} of strings",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{"a", "b", "c"},
				},
			},
			key:           "key",
			expectedSlice: []string{"a", "b", "c"},
			expectedOk:    true,
		},
		{
			name: "Value is []interface{} with non-string elements",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{"a", 2, "c"},
				},
			},
			key:           "key",
			expectedSlice: []string{},
			expectedOk:    false,
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"key": nil,
				},
			},
			key:           "key",
			expectedSlice: []string{},
			expectedOk:    false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:           "missing_key",
			expectedSlice: []string{},
			expectedOk:    false,
		},
		{
			name: "Value is of unexpected type (int)",
			params: &Params{
				Values: map[string]interface{}{
					"key": 123,
				},
			},
			key:           "key",
			expectedSlice: []string{},
			expectedOk:    false,
		},
		{
			name: "Value is empty string",
			params: &Params{
				Values: map[string]interface{}{
					"key": "",
				},
			},
			key:           "key",
			expectedSlice: []string{""},
			expectedOk:    true,
		},
		{
			name: "Value is empty []string",
			params: &Params{
				Values: map[string]interface{}{
					"key": []string{},
				},
			},
			key:           "key",
			expectedSlice: []string{},
			expectedOk:    true,
		},
		{
			name: "Value is empty []byte",
			params: &Params{
				Values: map[string]interface{}{
					"key": []byte(""),
				},
			},
			key:           "key",
			expectedSlice: []string{""},
			expectedOk:    true,
		},
		{
			name: "Value is string without commas",
			params: &Params{
				Values: map[string]interface{}{
					"key": "abc",
				},
			},
			key:           "key",
			expectedSlice: []string{"abc"},
			expectedOk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice, ok := tt.params.GetStringSliceOk(tt.key)
			assert.Equal(t, tt.expectedOk, ok, "Expected ok to be %v, got %v", tt.expectedOk, ok)
			assert.Equal(t, tt.expectedSlice, slice, "Expected slice to be %v, got %v", tt.expectedSlice, slice)
		})
	}
}

// TestParams_GetTimeInLocationOk tests the GetTimeInLocationOk method
func TestParams_GetTimeInLocationOk(t *testing.T) {
	loc := time.UTC // You can specify any location you need

	tests := []struct {
		name         string
		params       *Params
		key          string
		expectedTime time.Time
		expectedOk   bool
	}{
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:          "missing_key",
			expectedTime: time.Time{},
			expectedOk:   false,
		},
		{
			name: "Value is time.Time",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
			expectedOk:   true,
		},
		{
			name: "Value is string in RFC3339 format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15T12:30:45Z",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
			expectedOk:   true,
		},
		{
			name: "Value is string in DateOnly format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 0, 0, 0, 0, loc),
			expectedOk:   true,
		},
		{
			name: "Value is string in DateTime format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15 12:30:45",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
			expectedOk:   true,
		},
		{
			name: "Value is string in HTMLDateTimeLocal format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15T12:30",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 0, 0, loc),
			expectedOk:   true,
		},
		{
			name: "Value is string in invalid format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "invalid time format",
				},
			},
			key:          "timeKey",
			expectedTime: time.Time{},
			expectedOk:   false, // According to the function's behavior
		},
		{
			name: "Value is of unexpected type (int)",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": 1234567890,
				},
			},
			key:          "timeKey",
			expectedTime: time.Time{},
			expectedOk:   false, // According to the function's behavior
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": nil,
				},
			},
			key:          "timeKey",
			expectedTime: time.Time{},
			expectedOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeValue, ok := tt.params.GetTimeInLocationOk(tt.key, loc)
			assert.Equal(t, tt.expectedOk, ok, "Expected ok to be %v, got %v", tt.expectedOk, ok)
			assert.Equal(t, tt.expectedTime, timeValue, "Expected time to be %v, got %v", tt.expectedTime, timeValue)
		})
	}
}

// TestParams_GetTimeInLocation tests the GetTimeInLocation method
func TestParams_GetTimeInLocation(t *testing.T) {

	loc := time.UTC // You can specify any location you need

	tests := []struct {
		name         string
		params       *Params
		key          string
		expectedTime time.Time
	}{
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:          "missing_key",
			expectedTime: time.Time{},
		},
		{
			name: "Value is time.Time",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
		},
		{
			name: "Value is string in RFC3339 format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15T12:30:45Z",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 45, 0, time.UTC),
		},
		{
			name: "Value is string in DateOnly format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 0, 0, 0, 0, loc),
		},
		{
			name: "Value is string in DateTime format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15 12:30:45",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 45, 0, loc),
		},
		{
			name: "Value is string in HTMLDateTimeLocal format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "2023-10-15T12:30",
				},
			},
			key:          "timeKey",
			expectedTime: time.Date(2023, 10, 15, 12, 30, 0, 0, loc),
		},
		{
			name: "Value is string in invalid format",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": "invalid time format",
				},
			},
			key:          "timeKey",
			expectedTime: time.Time{},
		},
		{
			name: "Value is of unexpected type (int)",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": 1234567890,
				},
			},
			key:          "timeKey",
			expectedTime: time.Time{},
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"timeKey": nil,
				},
			},
			key:          "timeKey",
			expectedTime: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeValue := tt.params.GetTimeInLocation(tt.key, loc)
			assert.Equal(t, tt.expectedTime, timeValue, "Expected time to be %v, got %v", tt.expectedTime, timeValue)
		})
	}
}

// TestParams_GetJSONOk tests the GetJSONOk method
func TestParams_GetJSONOk(t *testing.T) {
	tests := []struct {
		name         string
		params       *Params
		key          string
		expectedData map[string]interface{}
		expectedOk   bool
	}{
		{
			name: "Value is map[string]interface{}",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": map[string]interface{}{
						"name": "John",
						"age":  30, // Stored as int
					},
				},
			},
			key: "jsonKey",
			expectedData: map[string]interface{}{
				"name": "John",
				"age":  30, // Expected as int
			},
			expectedOk: true,
		},
		{
			name: "Value is valid JSON string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `{"name":"Jane","age":25}`,
				},
			},
			key: "jsonKey",
			expectedData: map[string]interface{}{
				"name": "Jane",
				"age":  float64(25),
			},
			expectedOk: true,
		},
		{
			name: "Value is invalid JSON string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `{"name":"Invalid JSON", "age":}`,
				},
			},
			key:          "jsonKey",
			expectedData: nil,
			expectedOk:   false,
		},
		{
			name: "Value cannot be converted to string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": []string{"not", "a", "string"},
				},
			},
			key:          "jsonKey",
			expectedData: nil,
			expectedOk:   false,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:          "missingKey",
			expectedData: nil,
			expectedOk:   false,
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": nil,
				},
			},
			key:          "jsonKey",
			expectedData: nil,
			expectedOk:   false,
		},
		{
			name: "Value is of unexpected type (int)",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": 123,
				},
			},
			key:          "jsonKey",
			expectedData: nil,
			expectedOk:   false,
		},
		{
			name: "Value is empty string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": "",
				},
			},
			key:          "jsonKey",
			expectedData: nil,
			expectedOk:   false,
		},
		{
			name: "Value is JSON string representing array",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `["apple", "banana", "cherry"]`,
				},
			},
			key:          "jsonKey",
			expectedData: nil, // Function expects object, not array
			expectedOk:   false,
		},
		{
			name: "Value is valid JSON string with nested objects",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `{"person": {"name": "Alice", "age": 28}, "city": "Wonderland"}`,
				},
			},
			key: "jsonKey",
			expectedData: map[string]interface{}{
				"person": map[string]interface{}{
					"name": "Alice",
					"age":  float64(28),
				},
				"city": "Wonderland",
			},
			expectedOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, ok := tt.params.GetJSONOk(tt.key)
			assert.Equal(t, tt.expectedOk, ok, "Expected ok to be %v, got %v", tt.expectedOk, ok)
			assert.Equal(t, tt.expectedData, data, "Expected data to be %v, got %v", tt.expectedData, data)
		})
	}
}

// TestParams_GetJSON tests the GetJSON method
func TestParams_GetJSON(t *testing.T) {
	tests := []struct {
		name         string
		params       *Params
		key          string
		expectedData map[string]interface{}
	}{
		{
			name: "Value is map[string]interface{}",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": map[string]interface{}{
						"name": "John",
						"age":  30, // Stored as int
					},
				},
			},
			key: "jsonKey",
			expectedData: map[string]interface{}{
				"name": "John",
				"age":  30, // Expected as int
			},
		},
		{
			name: "Value is valid JSON string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `{"name":"Jane","age":25}`,
				},
			},
			key: "jsonKey",
			expectedData: map[string]interface{}{
				"name": "Jane",
				"age":  float64(25), // Decoded from JSON string
			},
		},
		{
			name: "Value is invalid JSON string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `{"name":"Invalid JSON", "age":}`,
				},
			},
			key:          "jsonKey",
			expectedData: nil, // Parsing fails, GetJSON returns nil
		},
		{
			name: "Value cannot be converted to string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": []string{"not", "a", "string"},
				},
			},
			key:          "jsonKey",
			expectedData: nil, // GetStringOk fails, GetJSON returns nil
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:          "missingKey",
			expectedData: nil,
		},
		{
			name: "Value is nil",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": nil,
				},
			},
			key:          "jsonKey",
			expectedData: nil,
		},
		{
			name: "Value is of unexpected type (int)",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": 123,
				},
			},
			key:          "jsonKey",
			expectedData: nil,
		},
		{
			name: "Value is empty string",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": "",
				},
			},
			key:          "jsonKey",
			expectedData: nil, // Empty string cannot be parsed as JSON
		},
		{
			name: "Value is JSON string representing array",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `["apple", "banana", "cherry"]`,
				},
			},
			key:          "jsonKey",
			expectedData: nil, // Expected map[string]interface{}, but JSON is an array
		},
		{
			name: "Value is valid JSON string with nested objects",
			params: &Params{
				Values: map[string]interface{}{
					"jsonKey": `{"person": {"name": "Alice", "age": 28}, "city": "Wonderland"}`,
				},
			},
			key: "jsonKey",
			expectedData: map[string]interface{}{
				"person": map[string]interface{}{
					"name": "Alice",
					"age":  float64(28), // JSON numbers are float64
				},
				"city": "Wonderland",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.params.GetJSON(tt.key)
			assert.Equal(t, tt.expectedData, data, "Expected data to be %v, got %v", tt.expectedData, data)
		})
	}
}

// TestParams_GetUint64SliceOk tests the GetUint64SliceOk method
func TestParams_GetUint64SliceOk(t *testing.T) {
	tests := []struct {
		name         string
		params       *Params
		key          string
		expectedData []uint64
		expectedOk   bool
	}{
		{
			name: "Valid positive integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{1, 2, 3},
				},
			},
			key:          "key",
			expectedData: []uint64{1, 2, 3},
			expectedOk:   true,
		},
		{
			name: "Contains negative integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{1, -2, 3},
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
			expectedOk:   false,
		},
		{
			name: "Empty slice",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{},
				},
			},
			key:          "key",
			expectedData: []uint64{},
			expectedOk:   true,
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:          "missing_key",
			expectedData: []uint64{},
			expectedOk:   false,
		},
		{
			name: "Zero values",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{0, 1, 2},
				},
			},
			key:          "key",
			expectedData: []uint64{0, 1, 2},
			expectedOk:   true,
		},
		{
			name: "All negative integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{-1, -2, -3},
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
			expectedOk:   false,
		},
		{
			name: "Large integers within uint64 range",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{2147483647}, // Max int32
				},
			},
			key:          "key",
			expectedData: []uint64{2147483647},
			expectedOk:   true,
		},
		{
			name: "Value is string of comma-separated numbers",
			params: &Params{
				Values: map[string]interface{}{
					"key": "4,5,6",
				},
			},
			key:          "key",
			expectedData: []uint64{4, 5, 6},
			expectedOk:   true,
		},
		{
			name: "Value is string with negative numbers",
			params: &Params{
				Values: map[string]interface{}{
					"key": "7,-8,9",
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
			expectedOk:   false,
		},
		{
			name: "Value is []interface{} with integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{10, 11, 12},
				},
			},
			key:          "key",
			expectedData: []uint64{10, 11, 12},
			expectedOk:   true,
		},
		{
			name: "Value is []interface{} with strings",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{"13", "14", "15"},
				},
			},
			key:          "key",
			expectedData: []uint64{13, 14, 15},
			expectedOk:   true,
		},
		{
			name: "Value is []interface{} with negative numbers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{16, -17, 18},
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
			expectedOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, ok := tt.params.GetUint64SliceOk(tt.key)
			assert.Equal(t, tt.expectedOk, ok, "Expected ok to be %v, got %v", tt.expectedOk, ok)
			assert.Equal(t, tt.expectedData, data, "Expected data to be %v, got %v", tt.expectedData, data)
		})
	}
}

// TestParams_GetUint64Slice tests the GetUint64Slice method
func TestParams_GetUint64Slice(t *testing.T) {
	tests := []struct {
		name         string
		params       *Params
		key          string
		expectedData []uint64
	}{
		{
			name: "Valid positive integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{1, 2, 3},
				},
			},
			key:          "key",
			expectedData: []uint64{1, 2, 3},
		},
		{
			name: "Contains negative integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{1, -2, 3},
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
		},
		{
			name: "Empty slice",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{},
				},
			},
			key:          "key",
			expectedData: []uint64{},
		},
		{
			name: "Key does not exist",
			params: &Params{
				Values: map[string]interface{}{},
			},
			key:          "missing_key",
			expectedData: []uint64{},
		},
		{
			name: "Zero values",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{0, 1, 2},
				},
			},
			key:          "key",
			expectedData: []uint64{0, 1, 2},
		},
		{
			name: "All negative integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{-1, -2, -3},
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
		},
		{
			name: "Large integers within uint64 range",
			params: &Params{
				Values: map[string]interface{}{
					"key": []int{2147483647}, // Max int32
				},
			},
			key:          "key",
			expectedData: []uint64{2147483647},
		},
		{
			name: "Value is string of comma-separated numbers",
			params: &Params{
				Values: map[string]interface{}{
					"key": "4,5,6",
				},
			},
			key:          "key",
			expectedData: []uint64{4, 5, 6},
		},
		{
			name: "Value is string with negative numbers",
			params: &Params{
				Values: map[string]interface{}{
					"key": "7,-8,9",
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
		},
		{
			name: "Value is []interface{} with integers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{10, 11, 12},
				},
			},
			key:          "key",
			expectedData: []uint64{10, 11, 12},
		},
		{
			name: "Value is []interface{} with strings",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{"13", "14", "15"},
				},
			},
			key:          "key",
			expectedData: []uint64{13, 14, 15},
		},
		{
			name: "Value is []interface{} with negative numbers",
			params: &Params{
				Values: map[string]interface{}{
					"key": []interface{}{16, -17, 18},
				},
			},
			key:          "key",
			expectedData: []uint64{}, // Expected empty slice on failure
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.params.GetUint64Slice(tt.key)
			assert.Equal(t, tt.expectedData, data, "Expected data to be %v, got %v", tt.expectedData, data)
		})
	}
}

// TestParams_Permit tests the Permit method
func TestParams_Permit(t *testing.T) {
	tests := []struct {
		name         string
		initialVals  map[string]interface{}
		allowedKeys  []string
		expectedVals map[string]interface{}
	}{
		{
			name: "All keys are allowed",
			initialVals: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			allowedKeys: []string{"key1", "key2"},
			expectedVals: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "No keys are allowed",
			initialVals: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			allowedKeys:  []string{},
			expectedVals: map[string]interface{}{},
		},
		{
			name: "Some keys are allowed",
			initialVals: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			allowedKeys: []string{"key1", "key3"},
			expectedVals: map[string]interface{}{
				"key1": "value1",
				"key3": "value3",
			},
		},
		{
			name: "Allowed keys list is empty",
			initialVals: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			allowedKeys:  []string{},
			expectedVals: map[string]interface{}{},
		},
		{
			name:         "Params.Values is empty",
			initialVals:  map[string]interface{}{},
			allowedKeys:  []string{"key1", "key2"},
			expectedVals: map[string]interface{}{},
		},
		{
			name: "Allowed keys contain keys not present in Params.Values",
			initialVals: map[string]interface{}{
				"key1": "value1",
				"key3": "value3",
			},
			allowedKeys: []string{"key1", "key2"},
			expectedVals: map[string]interface{}{
				"key1": "value1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{
				Values: tt.initialVals,
			}

			params.Permit(tt.allowedKeys)

			assert.Equal(t, tt.expectedVals, params.Values, "Expected Params.Values to be %v, got %v", tt.expectedVals, params.Values)
		})
	}
}

// TestMakeHTTPRouterParsedReq tests the MakeHTTPRouterParsedReq function
func TestMakeHTTPRouterParsedReq(t *testing.T) {
	tests := []struct {
		name           string
		params         httprouter.Params
		expectedValues map[string]interface{}
	}{
		{
			name: "Param key contains 'id', value is uint64",
			params: httprouter.Params{
				httprouter.Param{Key: "user_id", Value: "12345"},
			},
			expectedValues: map[string]interface{}{
				"user_id": uint64(12345),
			},
		},
		{
			name: "Param key contains 'id', value is not uint64",
			params: httprouter.Params{
				httprouter.Param{Key: "user_id", Value: "not_a_number"},
			},
			expectedValues: map[string]interface{}{
				"user_id": "not_a_number",
			},
		},
		{
			name: "Param key does not contain 'id'",
			params: httprouter.Params{
				httprouter.Param{Key: "name", Value: "Alice"},
			},
			expectedValues: map[string]interface{}{
				"name": "Alice",
			},
		},
		{
			name: "Multiple params",
			params: httprouter.Params{
				httprouter.Param{Key: "user_id", Value: "12345"},
				httprouter.Param{Key: "session_id", Value: "fake_session_id"},
				httprouter.Param{Key: "action", Value: "login"},
				httprouter.Param{Key: "invalid_id", Value: "not_a_number"},
			},
			expectedValues: map[string]interface{}{
				"user_id":    uint64(12345),
				"session_id": "fake_session_id", // Cannot be parsed as uint64
				"action":     "login",
				"invalid_id": "not_a_number", // Cannot be parsed as uint64
			},
		},
		{
			name: "Param key contains 'id', value is empty string",
			params: httprouter.Params{
				httprouter.Param{Key: "product_id", Value: ""},
			},
			expectedValues: map[string]interface{}{
				"product_id": "",
			},
		},
		{
			name: "Param key contains 'id', value is negative number",
			params: httprouter.Params{
				httprouter.Param{Key: "item_id", Value: "-1"},
			},
			expectedValues: map[string]interface{}{
				"item_id": "-1",
			},
		},
		{
			name: "Param key does not contain 'id', value is numeric",
			params: httprouter.Params{
				httprouter.Param{Key: "count", Value: "42"},
			},
			expectedValues: map[string]interface{}{
				"count": "42",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handler to check the params
			handler := func(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				params := GetParams(r)
				assert.Equal(t, tt.expectedValues, params.Values)
			}

			// Wrap the handler
			wrappedHandler := MakeHTTPRouterParsedReq(handler)

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)

			// Create a response recorder
			rw := httptest.NewRecorder()

			// Call the wrapped handler
			wrappedHandler(rw, req, tt.params)
		})
	}
}
