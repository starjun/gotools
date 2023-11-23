package gotools

import (
	"encoding/json"
	"log"
)

// JsonEncode map/sturct转json字符串
// v  interface{}
// return: json字符串
func JsonEncode(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	return string(s)
}

// JsonDecode 解析json字符到interface{}
// v  interface{}
// return: error
func JsonDecode(data string, v interface{}) error {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		log.Println(err)
	}
	return err
}
