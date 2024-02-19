package natsutils

import "encoding/json"

// ToPayload converts an object to a json byte array if his type is not string or []byte
func ToPayload(obj any) []byte {
	switch obj.(type) {
	case string:
		return []byte(obj.(string))
	case []byte:
		return obj.([]byte)
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return []byte(err.Error())
	}
	return bytes
}
