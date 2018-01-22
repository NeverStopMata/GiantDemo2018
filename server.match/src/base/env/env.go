package env

import (
	"base/glog"
	"encoding/json"
	"io/ioutil"
)

var configData map[string]map[string]string

func Load(path string) bool {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Error("[配置] 读取失败 ", path, ",", err)
		return false
	}
	err = json.Unmarshal(file, &configData)
	if err != nil {
		glog.Error("[配置] 解析失败 ", path, ",", err)
		return false
	}
	return true
}

func Get(table, key string) string {
	t, ok := configData[table]
	if !ok {
		return ""
	}
	val, ok := t[key]
	if !ok {
		return ""
	}
	return val
}

func Global(key string) string {
	val, ok := configData["global"][key]
	if !ok {
		return ""
	}
	return val
}

func Set(key1, key2, val string) {
	if kmap, ok := configData[key1]; ok {
		kmap[key2] = val
	}
}

func User(key string) string {
	val, ok := configData["user"][key]
	if !ok {
		return ""
	}
	return val
}

func Room(key string) string {
	val, ok := configData["room"][key]
	if !ok {
		return ""
	}
	return val
}

func Video(key string) string {
	val, ok := configData["video"][key]
	if !ok {
		return ""
	}
	return val
}
