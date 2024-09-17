package parameters

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const testJSONParam = `{ "test": true }`

// TestGetParams_ParseJSONBody
func TestGetParams_ParseJSONBody(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), "POST", "test", strings.NewReader(testJSONParam))
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
	r, err := http.NewRequestWithContext(context.Background(), "POST", "test", strings.NewReader(testJSONParam))
	assert.NoError(b, err)
	r.Header.Set("Content-Type", "application/json")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	for i := 0; i < b.N; i++ {
		_ = GetParams(r)
	}
}

// TestGetParams_ParseJSONBodyContentType
func TestGetParams_ParseJSONBodyContentType(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), "POST", "test", strings.NewReader(testJSONParam))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/json; charset=utf8")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// TestGetParams_ParseNestedJSONBody
func TestGetParams_ParseNestedJSONBody(t *testing.T) {
	body := "{ \"test\": true, \"coordinate\": { \"lat\": 50.505, \"lon\": 10.101 }}"
	r, err := http.NewRequestWithContext(context.Background(), "POST", "test", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "test?test=true", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "test?test=true", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	for i := 0; i < b.N; i++ {
		_ = GetParams(r)
	}
}

// TestParams_GetStringOk
func TestParams_GetStringOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=string", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=string", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetStringOk("test")
	}
}

// TestParams_GetBoolOk
func TestParams_GetBoolOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=true", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=true", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetBoolOk("test")
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test="+testString, strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test="+testString, strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=true", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_ = params.GetBool("test")
	}
}

// TestParams_GetFloatOk
func TestParams_GetFloatOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123.1234", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123.1234", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetFloatOk("test")
	}
}

// TestParams_GetFloatOk2
func TestParams_GetFloatOk2(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=null", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetFloatOk("test")
	assert.Equal(t, float64(0), val)
	assert.Equal(t, true, ok)
}

// TestParams_GetIntOk
func TestParams_GetIntOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetIntOk("test")
	}
}

// TestParams_GetInt8Ok
func TestParams_GetInt8Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt8Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int8(123), val)
}

// TestParams_GetInt8TooSmall
func TestParams_GetInt8TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=-300", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt8("test")
	assert.Equal(t, int8(0), val)
}

// TestParams_GetInt8TooBig
func TestParams_GetInt8TooBig(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=300", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt8("test")
	assert.Equal(t, int8(0), val)
}

// TestParams_GetInt16Ok
func TestParams_GetInt16Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt16Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int16(123), val)

	val = params.GetInt16("test")
	assert.Equal(t, int16(123), val)
}

// TestParams_GetInt16TooSmall
func TestParams_GetInt16TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=-32769", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt16("test")
	assert.Equal(t, int16(0), val)
}

// TestParams_GetInt16TooBig
func TestParams_GetInt16TooBig(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=32769", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt16("test")
	assert.Equal(t, int16(0), val)
}

// TestParams_GetInt32Ok
func TestParams_GetInt32Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetInt32Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, int32(123), val)

	val = params.GetInt32("test")
	assert.Equal(t, int32(123), val)
}

// TestParams_GetInt32TooSmall
func TestParams_GetInt32TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/test?test=%d", math.MinInt32-1), strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt32("test")
	assert.Equal(t, int32(0), val)
}

// TestParams_GetInt32TooBig
func TestParams_GetInt32TooBig(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/test?test=%d", math.MaxInt32+1), strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt32("test")
	assert.Equal(t, int32(0), val)
}

// TestParams_GetInt64Ok
func TestParams_GetInt64Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
	assert.NoError(b, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	for i := 0; i < b.N; i++ {
		_, _ = params.GetInt64Ok("test")
	}
}

// TestParams_GetInt64TooSmall
func TestParams_GetInt64TooSmall(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/test?test=%d", 0), strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val := params.GetInt64("test")
	assert.Equal(t, int64(0), val)
}

// TestParams_GetUint64Ok
func TestParams_GetUint64Ok(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetUint64Ok("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, uint64(123), val)
}

// TestGetParams_Post
func TestGetParams_Post(t *testing.T) {
	body := "test=true"
	r, err := http.NewRequestWithContext(context.Background(), "POST", "test", strings.NewReader(body))
	assert.NoError(t, err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, present := params.Get("test")
	assert.Equal(t, true, present)
	assert.Equal(t, true, val)
}

// TestParams_GetTimeOk
func TestParams_GetTimeOk(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=2020-12-31", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetTimeOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, "2020-12-31 00:00:00 +0000 UTC", val.String())
}

// TestGetParams_Put
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

// TestGetParams_ParsePostUrlJSON
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

// TestGetParams_ParseJSONBodyMux
func TestGetParams_ParseJSONBodyMux(t *testing.T) {

	r, err := http.NewRequestWithContext(context.Background(), "POST", "/test/42", strings.NewReader(testJSONParam))
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

// TestImbue
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

// TestImbue_Time
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

// TestHasAll
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "test?test=true", strings.NewReader(""))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "test", strings.NewReader(body))
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
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=123", strings.NewReader(body))
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
			"ints": overflowValue,
		},
	}

	slice, ok := params.GetIntSliceOk("ints")

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
