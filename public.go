package main

import "encoding/json"

func Response2Dict(b []byte) (v map[string]interface{}) {
	_ = json.Unmarshal(b, &v)
	return v
}
