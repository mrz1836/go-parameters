package parameters

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
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

// TestGetParams
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

// TestParams_GetBoolInt
func TestParams_GetBoolInt(t *testing.T) {
	body := ""
	r, err := http.NewRequestWithContext(context.Background(), "GET", "/test?test=1", strings.NewReader(body))
	assert.NoError(t, err)

	r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))

	params := GetParams(r)

	val, ok := params.GetBoolOk("test")
	assert.Equal(t, true, ok)
	assert.Equal(t, true, val)
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
	r, err := http.NewRequestWithContext(context.Background(), "PUT", "test", strings.NewReader(body))
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
