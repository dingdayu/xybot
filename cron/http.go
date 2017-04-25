package cron

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var cookieJarList map[string]*cookiejar.Jar
var baseUrl = "https://wx2.qq.com/cgi-bin/mmwebwx-bin"

type Http struct {
	Client *http.Client
	Jar    *cookiejar.Jar
}

func init() {
	cookieJarList = make(map[string]*cookiejar.Jar)
}

//func (j cookiejar.Jar)toJson(ur string)   {
//	urls,_ := url.Parse(ur);
//	cookies := j.Cookies(urls);
//
//	jsn,_ := json.Marshal(cookies);
//	return string(jsn)
//}
//
//func (j cookiejar.Jar)setCookie(ur string, json string)  {
//	var ps []*http.Cookie
//	json.Unmarshal([]byte(json), &ps)
//	j.SetCookies(ur, ps)
//}

func NewHttp(uuid string) *Http {
	var cookieJar *cookiejar.Jar
	if uuid != "" {
		if tmp, ok := cookieJarList[uuid]; ok {
			cookieJar = tmp
		} else {

			cookieJar, _ = cookiejar.New(nil)
			cookieJarList[uuid] = cookieJar
		}
	} else {
		cookieJar, _ = cookiejar.New(nil)
	}

	client := &http.Client{
		Jar: cookieJar,
	}
	return &Http{Client: client, Jar: cookieJar}
}

func (h *Http) Get(url string, params map[string]string) string {

	if !strings.Contains(url, "?") {
		url += "?"
	}
	if len(params) > 0 {
		for k, v := range params {
			url += k + "=" + v + "&"
		}
		url = strings.TrimSuffix(url, "&")
	}

	resp, err := h.Client.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return string(body)
}

func (h *Http) PostMap(url string, params map[string]string) string {
	var post string
	for k, v := range params {
		post += k + "=" + v + "&"
	}
	url = strings.TrimSuffix(url, "&")

	return h.Post(url, post)
}

func (h *Http) Post(url string, post string) string {
	resp, err := h.Client.Post(url,
		"application/json, text/plain, */*",
		strings.NewReader(post))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(err.Error())
	}
	return string(body)
}

func (h *Http) PostForm(url string, data url.Values) string {
	resp, err := h.Client.PostForm(url, data)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

func (h Http) SaveImage(url string, uuid string) {
	res, err := h.Client.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s", url, err)
		return
	}
	//按分辨率目录保存图片
	Dirname := "tmp/"
	if !isDirExist(Dirname) {
		os.Mkdir(Dirname, 0755)
		fmt.Printf("dir %s created\n", Dirname)
	}
	//根据URL文件名创建文件
	dst, err := os.Create(Dirname + uuid + ".png")
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s\n"+url, err)
		return
	}
	// 写入文件
	io.Copy(dst, res.Body)
}

func isDirExist(path string) bool {
	p, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return p.IsDir()
	}
}
