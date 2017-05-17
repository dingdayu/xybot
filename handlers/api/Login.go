package api

import (
	"encoding/json"
	"fmt"
	"github.com/dingdayu/wxbot/cron"
	"github.com/dingdayu/wxbot/model"
	"io"
	"net/http"
)

func HelloJson(w http.ResponseWriter, r *http.Request) {
	// 定义返回的结构体
	type jsonType struct {
		// 这里遵循大写字母开头方可被开放
		// 原因在于自定义结构体里面的对象，需要可以被json包访问到，
		// 而go规定只有大写开头的才能被包外部访问，而类型属于go语言的基本结构
		Name string
		age  int
	}

	// 实例化一个结构体
	hello := jsonType{Name: "dingdayu", age: 23}
	// map类型同样的使用方法
	//hello := make(map[string]string)
	// 这里不遵循大写字母开头的问题
	//hello["Name"] = "dingdayu"
	//hello["age"] = 23

	// 将结构体或类型转json字符串 除channel,complex和函数几种类型外，都可以转json
	// 注意  json.Marshal() 返回的是字节 需要转 string()
	if j, err := json.Marshal(hello); err != nil {
		fmt.Fprint(w, "json error")
	} else {
		// 返回json的类型头信息
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(j))
	}
}

type RetT struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 获取登陆uuid
func GetUUID(w http.ResponseWriter, r *http.Request) {
	uuid, err := cron.GetUuid()

	ret := RetT{}
	if err != nil {
		ret = RetT{Code: 301, Msg: err.Error()}
	} else {
		url := "https://login.weixin.qq.com/qrcode/" + uuid
		data := map[string]string{}
		data["url"] = url
		data["uuid"] = uuid
		ret = RetT{200, "success", data}
	}
	RetJson(ret, w)
}

// 发送文本消息
func SendText(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	uuid := r.FormValue("uuid")
	content := r.FormValue("content")

	ret := RetT{}
	if v, ok := cron.WxMap[uuid]; ok {
		msgId, err := v.SendTextMsg(username, content)
		if err != nil {
			ret = RetT{Code: 301, Msg: err.Error()}
		} else {
			data := map[string]string{}
			data["msgid"] = msgId
			ret = RetT{200, "success", data}
		}
	} else {
		ret = RetT{Code: 401, Msg: "uuid errror"}
	}
	RetJson(ret, w)
}

// 获取用户状态
func GetStatus(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	ret := RetT{}
	if v, ok := cron.WxMap[uuid]; ok {
		data := map[string]string{}
		data["status"] = v.Statue
		ret = RetT{200, "success", data}
	} else {
		u := model.GetUUIDTByUUID(uuid)
		if u.Id.Valid() {
			data := map[string]string{}
			data["status"] = u.Status
			ret = RetT{Code: 200, Msg: "success", Data: data}
		} else {
			ret = RetT{Code: 401, Msg: "uuid errror"}
		}

	}
	RetJson(ret, w)
}

func GetAllStatus(w http.ResponseWriter, r *http.Request) {
	ret := RetT{}
	data := make(map[string]interface{})
	data["count"] = len(cron.WxMap)

	type listT struct {
		UUID     string `json:"uuid,omitempty"`
		NickName string `json:"nickname,omitempty"`
	}
	list := []listT{}
	for uuid, item := range cron.WxMap {
		list = append(list, listT{UUID: uuid, NickName: item.LoginUser.NickName})
	}
	data["list"] = list
	ret = RetT{200, "success", data}
	RetJson(ret, w)
}

// 登出
func Logout(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	ret := RetT{}
	if v, ok := cron.WxMap[uuid]; ok {
		err := v.Logout()
		if err != nil {
			if _, ok := cron.WxMap[uuid]; ok {
				ret = RetT{Code: 302, Msg: "logout error!"}
			} else {
				ret = RetT{Code: 200, Msg: "success"}
			}
		} else {
			ret = RetT{Code: 200, Msg: "success"}
		}
	} else {
		ret = RetT{Code: 401, Msg: "uuid errror"}
	}
	RetJson(ret, w)
}

// 返回json
func RetJson(v interface{}, w http.ResponseWriter) {
	// 返回json的类型头信息
	w.Header().Set("Content-Type", "application/json")
	if j, err := json.Marshal(v); err != nil {
		io.WriteString(w, string("json error"))
	} else {
		io.WriteString(w, string(j))
	}
}
