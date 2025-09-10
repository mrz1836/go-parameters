/*
Package parameters parses json, msg pack, or multi-part form data into a parameters object
*/
package parameters

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/ugorji/go/codec"
)

// Constants for parameters package
const (
	// ParamsKeyName standard key name for parameter data
	ParamsKeyName paramKey = "params"

	// DateOnly is only the date
	DateOnly = "2006-01-02"

	// DateTime is not recommended, rather use time.RFC3339
	DateTime = "2006-01-02 15:04:05"

	// HTMLDateTimeLocal is the format used by the input type datetime-local
	HTMLDateTimeLocal = "2006-01-02T15:04"

	// MaxSafeInt is the maximum safe integer value
	MaxSafeInt = 1 << 53 // 9007199254740992
)

// Variables for parameters package
var (
	typeOfTime      = reflect.TypeOf(time.Time{})
	typeOfPtrToTime = reflect.PointerTo(typeOfTime)
)

// paramKey used for context.WithValue
type paramKey string

// Params is the parameter values
type Params struct {
	isBinary bool
	Values   map[string]interface{}
}

// CustomTypeHandler custom type handler
type CustomTypeHandler func(field *reflect.Value, value interface{}) error

// type CustomTypeHandler func(field *reflect.Value, value interface{})

// CustomTypeSetter is used when Imbue is called on an object to handle unknown types
var CustomTypeSetter CustomTypeHandler

// Get the param by key, return interface
func (p *Params) Get(key string) (val interface{}, ok bool) {
	keys := strings.Split(key, ".")
	root := p.Values
	count := len(keys)
	for i := 0; i < count; i++ {
		val, ok = root[keys[i]]
		if ok && i < count-1 {
			root = val.(map[string]interface{})
		}
	}
	return
}

// GetFloatOk get param by key, return float
func (p *Params) GetFloatOk(key string) (float64, bool) {
	val, ok := p.Get(key)
	if stringValue, stringOk := val.(string); stringOk {
		var err error
		val, err = strconv.ParseFloat(stringValue, 64)
		ok = err == nil
	}
	if ok && val != nil {
		return val.(float64), true
	}
	return 0, true
}

// GetFloat get param by key, return float
func (p *Params) GetFloat(key string) float64 {
	val, _ := p.GetFloatOk(key)
	return val
}

// GetFloatSliceOk get param by key, return slice of floats
func (p *Params) GetFloatSliceOk(key string) ([]float64, bool) {
	val, ok := p.Get(key)
	if ok {
		switch v := val.(type) {
		case []float64:
			return v, true
		case string:
			if v == "" {
				return []float64{}, true
			}
			raw := strings.Split(v, ",")
			slice := make([]float64, 0, len(raw))
			for _, k := range raw {
				if num, err := strconv.ParseFloat(strings.TrimSpace(k), 64); err == nil {
					// Reject NaN and Infinity values
					if math.IsNaN(num) || math.IsInf(num, 0) {
						return []float64{}, false
					}
					slice = append(slice, num)
				} else {
					return []float64{}, false
				}
			}
			return slice, true
		case []byte:
			vStr := string(v)
			if vStr == "" {
				return []float64{}, true
			}
			raw := strings.Split(vStr, ",")
			slice := make([]float64, 0, len(raw))
			for _, k := range raw {
				if num, err := strconv.ParseFloat(strings.TrimSpace(k), 64); err == nil {
					// Reject NaN and Infinity values
					if math.IsNaN(num) || math.IsInf(num, 0) {
						return []float64{}, false
					}
					slice = append(slice, num)
				} else {
					return []float64{}, false
				}
			}
			return slice, true
		case []interface{}:
			slice := make([]float64, 0, len(v))
			for _, k := range v {
				switch innerVal := k.(type) {
				case float64:
					// Reject NaN and Infinity values
					if math.IsNaN(innerVal) || math.IsInf(innerVal, 0) {
						return []float64{}, false
					}
					slice = append(slice, innerVal)
				case string:
					if num, err := strconv.ParseFloat(innerVal, 64); err == nil {
						// Reject NaN and Infinity values
						if math.IsNaN(num) || math.IsInf(num, 0) {
							return []float64{}, false
						}
						slice = append(slice, num)
					} else {
						return []float64{}, false
					}
				default:
					return []float64{}, false
				}
			}
			return slice, true
		default:
			return []float64{}, false
		}
	}
	return []float64{}, false
}

// GetFloatSlice get param by key, return slice of floats
func (p *Params) GetFloatSlice(key string) []float64 {
	val, _ := p.GetFloatSliceOk(key)
	return val
}

// GetBoolOk get param by key, return boolean
func (p *Params) GetBoolOk(key string) (bool, bool) {
	val, ok := p.Get(key)
	if ok {
		if b, isBool := val.(bool); isBool {
			return b, true
		} else if i, isInt := p.GetIntOk(key); isInt {
			if i == 0 {
				return false, true
			}

			return true, true
		}
	}
	return false, false
}

// GetBool get param by key, return boolean
func (p *Params) GetBool(key string) bool {
	val, _ := p.GetBoolOk(key)
	return val
}

// GetIntOk get param by key, return integer
func (p *Params) GetIntOk(key string) (int, bool) {
	val, ok := p.Get(key)
	if !ok || val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return v, true
	case int8, int16, int32, int64:
		i := reflect.ValueOf(v).Int()
		if i >= int64(math.MinInt) && i <= int64(math.MaxInt) {
			return int(i), true
		}
		return 0, false // Overflow
	case uint, uint8, uint16, uint32, uint64:
		u := reflect.ValueOf(v).Uint()
		if u <= uint64(math.MaxInt) {
			return int(u), true
		}
		return 0, false // Overflow
	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f >= float64(math.MinInt) && f <= float64(math.MaxInt) && f == math.Trunc(f) {
			if f > MaxSafeInt || f < -MaxSafeInt {
				return 0, false // Value exceeds safe integer range
			}
			return int(f), true
		}
		return 0, false // Overflow or non-integer float
	case string:
		if parsedInt, err := strconv.ParseInt(v, 10, 64); err == nil {
			if parsedInt >= int64(math.MinInt) && parsedInt <= int64(math.MaxInt) {
				return int(parsedInt), true
			}
			return 0, false // Overflow
		}
		return 0, false // Parsing failed
	case []byte:
		s := string(v)
		if parsedInt, err := strconv.ParseInt(s, 10, 64); err == nil {
			if parsedInt >= int64(math.MinInt) && parsedInt <= int64(math.MaxInt) {
				return int(parsedInt), true
			}
			return 0, false // Overflow
		}
		return 0, false // Parsing failed
	default:
		return 0, false
	}
}

// GetInt get param by key, return integer
func (p *Params) GetInt(key string) int {
	val, _ := p.GetIntOk(key)
	return val
}

// GetInt8Ok get param by key, return integer
func (p *Params) GetInt8Ok(key string) (int8, bool) {
	val, ok := p.GetIntOk(key)

	if !ok || val < math.MinInt8 || val > math.MaxInt8 {
		return 0, false
	}

	return int8(val), true
}

// GetInt8 get param by key, return integer
func (p *Params) GetInt8(key string) int8 {
	val, _ := p.GetInt8Ok(key)
	return val
}

// GetInt16Ok get param by key, return integer
func (p *Params) GetInt16Ok(key string) (int16, bool) {
	val, ok := p.GetIntOk(key)

	if !ok || val < math.MinInt16 || val > math.MaxInt16 {
		return 0, false
	}

	return int16(val), true
}

// GetInt16 get param by key, return integer
func (p *Params) GetInt16(key string) int16 {
	val, _ := p.GetInt16Ok(key)
	return val
}

// GetInt32Ok get param by key, return integer
func (p *Params) GetInt32Ok(key string) (int32, bool) {
	val, ok := p.GetIntOk(key)

	if !ok || val < math.MinInt32 || val > math.MaxInt32 {
		return 0, false
	}

	return int32(val), true
}

// GetInt32 get param by key, return integer
func (p *Params) GetInt32(key string) int32 {
	val, _ := p.GetInt32Ok(key)
	return val
}

// GetInt64Ok get param by key, return integer
func (p *Params) GetInt64Ok(key string) (int64, bool) {
	val, ok := p.Get(key)
	if !ok || val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case int64:
		return v, true
	case int, int8, int16, int32:
		return reflect.ValueOf(v).Int(), true
	case uint, uint8, uint16, uint32, uint64:
		u := reflect.ValueOf(v).Uint()
		if u <= uint64(math.MaxInt64) {
			return int64(u), true
		}
		return 0, false // Overflow
	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f >= float64(math.MinInt64) && f <= float64(math.MaxInt64) {
			// Check if the float is an integer value
			if f == math.Trunc(f) {
				return int64(f), true
			}
		}
		return 0, false // Overflow or non-integer float
	case string:
		if parsedInt, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsedInt, true
		}
		return 0, false // Parsing failed
	case []byte:
		s := string(v)
		if parsedInt, err := strconv.ParseInt(s, 10, 64); err == nil {
			return parsedInt, true
		}
		return 0, false // Parsing failed
	default:
		return 0, false
	}
}

// GetInt64 get param by key, return integer
func (p *Params) GetInt64(key string) int64 {
	val, _ := p.GetInt64Ok(key)
	return val
}

// GetIntSliceOk get param by key, return slice of integers
func (p *Params) GetIntSliceOk(key string) ([]int, bool) {
	val, ok := p.Get(key)
	if ok {
		switch v := val.(type) {
		case []int:
			return v, true
		case []byte:
			valStr := string(v)
			raw := strings.Split(valStr, ",")
			slice := make([]int, 0, len(raw))
			for _, k := range raw {
				if num, err := strconv.ParseInt(k, 10, 64); err == nil {
					if num >= int64(math.MinInt) && num <= int64(math.MaxInt) {
						slice = append(slice, int(num))
					} else {
						return slice, false
					}
				} else {
					return slice, false
				}
			}
			return slice, true
		case string:
			if len(v) > 0 {
				raw := strings.Split(v, ",")
				slice := make([]int, 0, len(raw))
				for _, k := range raw {
					if num, err := strconv.ParseInt(k, 10, 64); err == nil {
						if num >= int64(math.MinInt) && num <= int64(math.MaxInt) {
							slice = append(slice, int(num))
						} else {
							return slice, false
						}
					} else {
						return slice, false
					}
				}
				return slice, true
			}
			return nil, true
		case []interface{}:
			raw := v
			slice := make([]int, 0, len(raw))
			for _, k := range raw {
				switch num := k.(type) {
				case int:
					slice = append(slice, num)
				case float64:
					if num >= float64(math.MinInt) && num <= float64(math.MaxInt) {
						slice = append(slice, int(num))
					} else {
						return slice, false
					}
				case string:
					if parsed, err := strconv.ParseInt(num, 10, 64); err == nil {
						if parsed >= int64(math.MinInt) && parsed <= int64(math.MaxInt) {
							slice = append(slice, int(parsed))
						} else {
							return slice, false
						}
					} else {
						return slice, false
					}
				default:
					return slice, false
				}
			}
			return slice, true
		default:
			return nil, true
		}
	}
	return []int{}, false
}

// GetIntSlice get param by key, return slice of integers
func (p *Params) GetIntSlice(key string) []int {
	val, _ := p.GetIntSliceOk(key)
	return val
}

// GetUint64Ok get param by key, return unsigned integer
func (p *Params) GetUint64Ok(key string) (uint64, bool) {
	val, ok := p.Get(key)
	if !ok || val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case uint64:
		return v, true
	case uint, uint8, uint16, uint32:
		return reflect.ValueOf(v).Uint(), true
	case int, int8, int16, int32, int64:
		i := reflect.ValueOf(v).Int()
		if i >= 0 {
			return uint64(i), true
		}
		return 0, false
	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f >= 0 && f == math.Trunc(f) && f <= float64(math.MaxUint64) {
			return uint64(f), true
		}
		return 0, false
	case string:
		if parsedUint, err := strconv.ParseUint(v, 10, 64); err == nil {
			return parsedUint, true
		}
		// Do not parse strings as float64; return false if parsing fails
		return 0, false
	case []byte:
		s := string(v)
		if parsedUint, err := strconv.ParseUint(s, 10, 64); err == nil {
			return parsedUint, true
		}
		// Do not parse []byte as float64; return false if parsing fails
		return 0, false
	default:
		return 0, false
	}
}

// GetUint64 get param by key, return unsigned integer
func (p *Params) GetUint64(key string) uint64 {
	val, _ := p.GetUint64Ok(key)
	return val
}

// GetUint64SliceOk get param by key, return slice of unsigned integers
func (p *Params) GetUint64SliceOk(key string) ([]uint64, bool) {
	if raw, ok := p.GetIntSliceOk(key); ok {
		slice := make([]uint64, len(raw))
		for i, num := range raw {
			if num < 0 {
				// Return empty slice instead of nil
				return []uint64{}, false
			}
			slice[i] = uint64(num)
		}
		return slice, true
	}

	return []uint64{}, false
}

// GetUint64Slice get param by key, return slice of unsigned integers
func (p *Params) GetUint64Slice(key string) []uint64 {
	val, _ := p.GetUint64SliceOk(key)
	return val
}

// GetStringOk get param by key, return string
func (p *Params) GetStringOk(key string) (string, bool) {
	val, ok := p.Get(key)
	if ok {
		if s, is := val.(string); is {
			return s, true
		} else if b, good := val.([]byte); good {
			return string(b), true
		}
	}
	return "", false
}

// GetString get param by key, return string
func (p *Params) GetString(key string) string {
	val, _ := p.GetStringOk(key)
	return strings.Trim(val, " ")
}

// GetStringSliceOk get param by key, return slice of strings
func (p *Params) GetStringSliceOk(key string) ([]string, bool) {
	val, ok := p.Get(key)
	if ok {
		switch v := val.(type) {
		case []string:
			return v, true
		case []byte:
			return strings.Split(string(v), ","), true
		case string:
			return strings.Split(v, ","), true
		case []interface{}:
			slice := make([]string, 0, len(v))
			for _, k := range v {
				if str, okS := k.(string); okS {
					slice = append(slice, str)
				} else {
					return []string{}, false
				}
			}
			return slice, true
		default:
			return []string{}, false
		}
	}
	return []string{}, false
}

// GetStringSlice get param by key, return slice of strings
func (p *Params) GetStringSlice(key string) []string {
	val, _ := p.GetStringSliceOk(key)
	return val
}

// GetBytesOk get param by key, return slice of bytes
func (p *Params) GetBytesOk(key string) ([]byte, bool) {
	if dataStr, ok := p.Get(key); ok {
		var dataByte []byte
		if dataByte, ok = dataStr.([]byte); !ok {
			var err error
			dataByte, err = base64.StdEncoding.DecodeString(dataStr.(string))
			if err != nil {
				log.Println("error decoding data:", key, err)
				return nil, true
			}
			p.Values[key] = dataByte
		}
		return dataByte, true
	}
	return nil, false
}

// GetBytes get param by key, return slice of bytes
func (p *Params) GetBytes(key string) []byte {
	val, _ := p.GetBytesOk(key)
	return val
}

// GetTimeOk get param by key, return time
func (p *Params) GetTimeOk(key string) (time.Time, bool) {
	return p.GetTimeInLocationOk(key, time.UTC)
}

// GetTime get param by key, return time
func (p *Params) GetTime(key string) time.Time {
	val, _ := p.GetTimeOk(key)
	return val
}

// GetTimeInLocationOk get param by key, return time
func (p *Params) GetTimeInLocationOk(key string, loc *time.Location) (time.Time, bool) {
	val, ok := p.Get(key)
	if !ok {
		return time.Time{}, false
	}
	if t, success := val.(time.Time); success {
		return t, true
	}
	if str, success := val.(string); success {
		if t, err := time.ParseInLocation(time.RFC3339, str, loc); err == nil {
			return t, true
		}
		if t, err := time.ParseInLocation(DateOnly, str, loc); err == nil {
			return t, true
		}
		if t, err := time.ParseInLocation(DateTime, str, loc); err == nil {
			return t, true
		}
		if t, err := time.ParseInLocation(HTMLDateTimeLocal, str, loc); err == nil {
			return t, true
		}
	}
	return time.Time{}, false // Changed from true to false
}

// GetTimeInLocation get param by key, return time
func (p *Params) GetTimeInLocation(key string, loc *time.Location) time.Time {
	val, _ := p.GetTimeInLocationOk(key, loc)
	return val
}

// GetFileOk get param by key, return file
func (p *Params) GetFileOk(key string) (*multipart.FileHeader, bool) {
	val, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	if fh, good := val.(*multipart.FileHeader); good {
		return fh, true
	}
	return nil, true
}

// GetJSONOk get param by key, return map of string interface
func (p *Params) GetJSONOk(key string) (map[string]interface{}, bool) {
	if v, ok := p.Get(key); ok {
		if d, good := v.(map[string]interface{}); good {
			return d, true
		}
	}
	val, ok := p.GetStringOk(key)
	var jsonData map[string]interface{}
	if !ok {
		return jsonData, false
	}
	if err := json.NewDecoder(strings.NewReader(val)).Decode(&jsonData); err != nil {
		return jsonData, false
	}
	return jsonData, true
}

// GetJSON get param by key, return map of string interface
func (p *Params) GetJSON(key string) map[string]interface{} {
	val, _ := p.GetJSONOk(key)
	return val
}

// Clone makes a copy of this params object
func (p *Params) Clone() *Params {
	values := make(map[string]interface{}, len(p.Values))
	for k, v := range p.Values {
		values[k] = v
	}
	return &Params{
		isBinary: p.isBinary,
		Values:   values,
	}
}

// Imbue sets the parameters to the object by type; does not handle nested parameters
func (p *Params) Imbue(obj interface{}) {
	// Get the type of the object
	typeOfObject := reflect.TypeOf(obj).Elem()

	// Get the object
	objectValue := reflect.ValueOf(obj).Elem()

	// Loop our parameters
	for k := range p.Values {

		// Make the incoming key_name into KeyName
		key := SnakeToCamelCase(k, true)

		// Get the type and bool if found
		fieldType, found := typeOfObject.FieldByName(key)

		// Skip parameter if not found on struct
		if !found {
			continue
		}

		// Get the field of the key
		field := objectValue.FieldByName(key)

		// Check our types and set accordingly
		if fieldType.Type.Kind() == reflect.String {
			// Set string
			field.Set(reflect.ValueOf(p.GetString(k)))
		} else if fieldType.Type.Kind() == reflect.Uint64 {
			// Set Uint64
			field.Set(reflect.ValueOf(p.GetUint64(k)))
		} else if fieldType.Type.Kind() == reflect.Int {
			// Set Int
			field.Set(reflect.ValueOf(p.GetInt(k)))
		} else if fieldType.Type.Kind() == reflect.Bool {
			// Set bool
			field.Set(reflect.ValueOf(p.GetBool(k)))
		} else if fieldType.Type.Kind() == reflect.Float32 {
			// Set float32
			field.Set(reflect.ValueOf(float32(p.GetFloat(k))))
		} else if fieldType.Type.Kind() == reflect.Float64 {
			// Set float64
			field.Set(reflect.ValueOf(p.GetFloat(k)))
		} else if fieldType.Type == reflect.SliceOf(reflect.TypeOf("")) {
			// Set []string
			field.Set(reflect.ValueOf(p.GetStringSlice(k)))
		} else if fieldType.Type == reflect.SliceOf(reflect.TypeOf(0)) {
			// Set []int
			field.Set(reflect.ValueOf(p.GetIntSlice(k)))
		} else if fieldType.Type == reflect.SliceOf(reflect.TypeOf(uint64(0))) {
			// Set []uint64
			field.Set(reflect.ValueOf(p.GetUint64Slice(k)))
		} else if fieldType.Type == reflect.SliceOf(reflect.TypeOf(float64(0))) {
			// Set []float64
			field.Set(reflect.ValueOf(p.GetFloatSlice(k)))
		} else if fieldType.Type == typeOfTime {
			// Set time.Time
			field.Set(reflect.ValueOf(p.GetTime(k)))
		} else if fieldType.Type == typeOfPtrToTime {
			// Set *time.Time
			t := p.GetTime(k)
			field.Set(reflect.ValueOf(&t))
		} else {
			val, _ := p.Get(k)
			if CustomTypeSetter != nil && CustomTypeSetter(&field, val) == nil {
				continue
			}

			if subValues, ok := p.GetJSONOk(k); ok {
				fieldValue := reflect.Indirect(objectValue).FieldByName(key)
				if reflect.ValueOf(fieldValue).IsZero() {
					continue
				}

				typeOfP := reflect.TypeOf(fieldValue.Interface())
				newObj := reflect.New(typeOfP).Interface()

				subParam := &Params{
					Values: subValues,
				}
				subParam.Imbue(newObj)
				field.Set(reflect.ValueOf(newObj).Elem())
			}
		}
	}
}

// HasAll will return if all specified keys are found in the params object
func (p *Params) HasAll(keys ...string) (bool, []string) {
	missing := make([]string, 0)
	for _, key := range keys {
		if _, exists := p.Values[key]; !exists {
			missing = append(missing, key)
		}
	}
	return len(missing) == 0, missing
}

// Permit permits only the allowed fields given by allowedKeys
func (p *Params) Permit(allowedKeys []string) {
	for key := range p.Values {
		if !contains(allowedKeys, key) {
			delete(p.Values, key)
		}
	}
}

// contains contains needle in haystack
func contains(haystack []string, needle string) bool {
	needle = strings.ToLower(needle)
	for _, straw := range haystack {
		if strings.ToLower(straw) == needle {
			return true
		}
	}
	return false
}

// GetParams get parameters
func GetParams(req *http.Request) *Params {
	params, ok := req.Context().Value(ParamsKeyName).(*Params)
	if !ok {
		return nil
	}
	return params
}

// ParseParams parse parameters
func ParseParams(req *http.Request) *Params {
	var p Params
	if params, exists := req.Context().Value(ParamsKeyName).(*Params); exists {
		return params
	}
	ct := req.Header.Get("Content-Type")
	ct = strings.Split(ct, ";")[0]
	if ct == "multipart/form-data" {
		if err := req.ParseMultipartForm(10000000); err != nil {
			log.Println("Request.ParseMultipartForm error:", err)
		}
	} else {
		if err := req.ParseForm(); err != nil {
			log.Println("request.ParseForm error:", err)
		}
	}
	tempMap := make(map[string]interface{}, len(req.Form))
	for k, v := range req.Form {
		if strings.ToLower(v[0]) == "true" {
			tempMap[k] = true
		} else if strings.ToLower(v[0]) == "false" {
			tempMap[k] = false
		} else {
			tempMap[k] = v[0]
		}
	}

	if req.MultipartForm != nil {
		for k, v := range req.MultipartForm.File {
			tempMap[k] = v[0]
		}
	}

	// read the whole body into bytes
	body, err := io.ReadAll(req.Body)
	if err == nil {
		// must close
		if err = req.Body.Close(); err == nil {
			// no errors, restore the body on the request for other readers
			req.Body = io.NopCloser(bytes.NewReader(body))
		}
	}

	if ct == "application/json" && req.ContentLength > 0 {
		err = json.Unmarshal(body, &p.Values)
		if err != nil {
			log.Println("content-type is \"application/json\" but no valid json data received:", err)
			p.Values = tempMap
		}
		for k, v := range tempMap {
			if _, pres := p.Values[k]; !pres {
				p.Values[k] = v
			}
		}
	} else if ct == "application/x-msgpack" {
		var mh codec.MsgpackHandle
		p.isBinary = true
		mh.MapType = reflect.TypeOf(p.Values)
		if len(body) > 0 {
			buff := bytes.NewBuffer(body)
			first := body[0]
			if (first >= 0x80 && first <= 0x8f) || (first == 0xde || first == 0xdf) {
				err = codec.NewDecoder(buff, &mh).Decode(&p.Values)
				if err != nil && errors.Is(err, io.EOF) {
					log.Println("failed decoding msgpack:", err)
				}
			} else {
				if p.Values == nil {
					p.Values = make(map[string]interface{})
				}
				// var err error
				for err == nil {
					paramValues := make([]interface{}, 0)
					err = codec.NewDecoder(buff, &mh).Decode(&paramValues)
					if err != nil && errors.Is(err, io.EOF) {
						log.Println("failed decoding msgpack:", err)
					} else {
						for i := len(paramValues) - 1; i >= 1; i -= 2 {
							// Safely convert key to string, handling different types
							var keyStr string
							switch keyVal := paramValues[i-1].(type) {
							case []byte:
								keyStr = string(keyVal)
							case string:
								keyStr = keyVal
							case int64, int, int32, int16, int8:
								keyStr = fmt.Sprintf("%d", keyVal)
							case uint64, uint, uint32, uint16, uint8:
								keyStr = fmt.Sprintf("%d", keyVal)
							case float64, float32:
								keyStr = fmt.Sprintf("%.0f", keyVal)
							default:
								// Skip this key-value pair if key type is not supported
								continue
							}
							p.Values[keyStr] = paramValues[i]
						}
					}
				}
			}
		} else {
			p.Values = make(map[string]interface{})
		}
		for k, v := range tempMap {
			if _, pres := p.Values[k]; !pres {
				p.Values[k] = v
			}
		}
	} else {
		p.Values = tempMap
	}

	for k, v := range mux.Vars(req) {
		const keyID = "id"
		if strings.Contains(k, keyID) {
			var id uint64
			id, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				p.Values[k] = v
			} else {
				p.Values[k] = id
			}
		} else {
			p.Values[k] = v
		}
	}

	return &p
}

// MakeParsedReq make parsed request
func MakeParsedReq(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))
		fn(rw, r)
	}
}

// MakeHTTPRouterParsedReq make http router parsed request
func MakeHTTPRouterParsedReq(fn httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKeyName, ParseParams(r)))
		params := GetParams(r)
		for _, param := range p {
			const keyID = "id"
			if strings.Contains(param.Key, keyID) {
				id, err := strconv.ParseUint(param.Value, 10, 64)
				if err != nil {
					params.Values[param.Key] = param.Value
				} else {
					params.Values[param.Key] = id
				}
			} else {
				params.Values[param.Key] = param.Value
			}
		}
		fn(rw, r, p)
	}
}
