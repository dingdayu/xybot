package utils

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"regexp"
)

func PregMatch(pattern string, content string) (str []string) {
	reg := regexp.MustCompile(pattern)
	me := reg.FindAllStringSubmatch(content, -1)
	for _, v := range me {
		if len(v) > 2 {
			for i := 1; i < len(v); i++ {
				str = append(str, v[i])
			}
		} else {
			str = append(str, v[1])
		}
	}
	return str
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Struct2BsonMap(obj interface{}) bson.M {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(bson.M)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// 有问题，待调试
//func Map2Struct(m map[string]interface{}, s interface{})  {
//	t := reflect.TypeOf(s)
//	val := reflect.ValueOf(s)
//
//	for k,v := range m {
//		if k == t.FieldByName(k) {
//			val.FieldByName(k).Set(v)
//		}
//	}
//	return m
//}

func Struct2Struct(o interface{}, n interface{}) {
	js, _ := json.Marshal(o)
	json.Unmarshal(js, n)
}
