package api

import (
	"fmt"
	"github.com/dingdayu/wxbot/cron"
	"github.com/dingdayu/wxbot/utils"
	"io"
	"net/http"
	"os"
	"strconv"
)

func UploadHandle(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	username := r.FormValue("username")
	ret := RetT{}
	if uuid == "" {
		ret = RetT{Code: 302, Msg: "uuid Not Empty"}
		RetJson(ret, w)
		return
	}
	us := cron.WxLoginStatus{}
	if u, ok := cron.WxMap[uuid]; ok {
		us = *u
	} else {
		ret = RetT{Code: 401, Msg: "uuid errror"}
		RetJson(ret, w)
		return
	}

	if "POST" == r.Method {
		file, head, err := r.FormFile("file")
		if err != nil {
			ret = RetT{Code: 303, Msg: "File Error"}
			RetJson(ret, w)
			return
		}
		defer file.Close()
		//创建文件
		path := "./tmp/" + strconv.Itoa(us.LoginUser.Uin) + "/" + uuid + "/"

		if !utils.IsDirExist(path) {
			os.MkdirAll(path, 0755)
			fmt.Printf("dir %s created\n", path)
		}
		//根据URL文件名创建文件
		f, err := os.Create(path + head.Filename)
		defer f.Close()
		io.Copy(f, file)

		msgid, err := us.SendImagesMsg(username, path+head.Filename)
		if err == nil {
			data := map[string]string{}
			data["msgid"] = msgid
			ret = RetT{Code: 200, Msg: "success", Data: data}
		} else {
			ret = RetT{Code: 304, Msg: "MediaId Error:" + err.Error()}
		}

	} else {
		ret = RetT{Code: 301, Msg: "Method Error"}
	}
	RetJson(ret, w)
}
