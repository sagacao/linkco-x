package xutils

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	"github.com/nothollyhigh/kiss/log"
	"github.com/nothollyhigh/kiss/util"
)

// 定义JSON操作
var (
	json              = jsoniter.ConfigCompatibleWithStandardLibrary
	JSONMarshal       = json.Marshal
	JSONUnmarshal     = json.Unmarshal
	JSONMarshalIndent = json.MarshalIndent
	JSONNewDecoder    = json.NewDecoder
	JSONNewEncoder    = json.NewEncoder
)

// JSONMarshalToString JSON编码为字符串
func JSONMarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		log.Error("JSONMarshalToString error:[%v]", err)
		return ""
	}
	return s
}

// JSONParse json to GO TYPE
func JSONParse(data string, obj interface{}) error {
	err := jsoniter.UnmarshalFromString(data, &obj)
	if err != nil {
		log.Error("JSONParse Failed: %v", err)
		return err
	}
	return nil
}

// ReadJSON read json to obj
func ReadJSON(obj interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data = util.TrimComment(data)
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}
