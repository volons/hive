package libs

import (
	"encoding/json"
	"math"
)

// ParseJSON parses a JSON string and returns a JSONObject
func ParseJSON(jsonStr []byte) (JSONObject, error) {
	data := make(JSONObject)
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ToJSON encodes any data into a JSON string
func ToJSON(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ToJSONString encodes any data into a JSON string if it fails it returns null
func ToJSONString(data interface{}) string {
	str, err := ToJSON(data)
	if err != nil {
		return "null"
	}

	return str
}

// JSONData is the common interface for JSONObject and JSONBool
type JSONData interface {
	JSON() (string, error)
	String() string
	Type() string
}

// JSONObject provides convenience methods for maps returned
// by the go JSON parser
type JSONObject map[string]interface{}

// JSON encodes this JSONData object into a JSON string
func (data JSONObject) JSON() (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// String encodes this JSONData object into a JSON string
// Like JSON but returns "{}" if it fails
func (data JSONObject) String() string {
	str, err := data.JSON()
	if err != nil {
		return "{}"
	}

	return str
}

// HasKey returns true if the key is present
func (data JSONObject) HasKey(key string) bool {
	_, ok := data[key]
	return ok
}

// GetKeys returns an array of all the keys in this JSON object
func (data JSONObject) GetKeys() []string {
	keys := make([]string, len(data))
	i := 0
	for key := range data {
		keys[i] = key
		i++
	}

	return keys
}

// GetAny returns any value from the JSON object
func (data JSONObject) GetAny(key string) (interface{}, bool) {
	val, ok := data[key]
	return val, ok
}

// GetString returns a value from the JSON object as string
func (data JSONObject) GetString(key string) (string, bool) {
	val := data[key]

	if val != nil {
		out, ok := val.(string)
		return out, ok
	}

	return "", false
}

// GetObj returns a value from the JSON object as JsonData
func (data JSONObject) GetObj(key string) (JSONObject, bool) {
	val := data[key]

	if val != nil {
		out, ok := val.(map[string]interface{})
		return out, ok
	}

	return nil, false
}

// GetArray returns a value from the JSON object as JsonData
func (data JSONObject) GetArray(key string) (JSONArray, bool) {
	val := data[key]

	if val != nil {
		out, ok := val.([]interface{})
		return out, ok
	}

	return nil, false
}

// GetNumber returns a value from the JSON object as a float64
func (data JSONObject) GetNumber(key string) (float64, bool) {
	val := data[key]

	if val != nil {
		out, ok := val.(float64)
		return out, ok
	}

	return 0, false
}

// GetNum returns a value from the JSON object as a float64 or
// nil if not present or not a number
func (data JSONObject) GetNum(key string) *float64 {
	val := data[key]

	if val != nil {
		out, ok := val.(float64)
		if ok {
			return &out
		}
	}

	return nil
}

// GetBool returns a value from the JSON object as a boolean
func (data JSONObject) GetBool(key string) (bool, bool) {
	val := data[key]

	if val != nil {
		out, ok := val.(bool)
		return out, ok
	}

	return false, false
}

// SetNumber sets a float64 number on the object if it is not NaN
func (data JSONObject) SetNumber(key string, val float64) bool {
	if math.IsNaN(val) {
		return false
	}

	data[key] = val
	return true
}

// Type returns the type of this JSONData entity
func (data JSONObject) Type() string {
	return "object"
}

// JSONArray provides convenience methods for arrays returned
// by the go JSON parser
type JSONArray []interface{}

// JSON encodes this JSONData object into a JSON string
func (data JSONArray) JSON() (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// String encodes this JSONData object into a JSON string
// Like JSON but returns "{}" if it fails
func (data JSONArray) String() string {
	str, err := data.JSON()
	if err != nil {
		return "[]"
	}

	return str
}

// JSONBool provides json convenience methods for booleans
type JSONBool bool

// JSON encodes this boolean into a JSON string
func (b JSONBool) JSON() (string, error) {
	if b {
		return "true", nil
	}

	return "false", nil
}

// String encodes this boolean into a JSON string
func (b JSONBool) String() string {
	if b {
		return "true"
	}

	return "false"
}

// Type returns the type of this JSONData entity
func (b JSONBool) Type() string {
	return "boolean"
}
