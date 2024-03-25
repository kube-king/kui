package json_util

import "encoding/json"

func ToJsonString(data interface{}) string {
	marshal, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(marshal)
}
