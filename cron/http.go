package cron

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dingdayu/wxbot/utils"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var cookieJarList map[string]*cookiejar.Jar
var baseUrl = "https://wx2.qq.com/cgi-bin/mmwebwx-bin"
var mediaCount = 0

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

// get发送字符串的map调用
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

// post发送字符串的map调用
func (h *Http) PostMap(url string, params map[string]string) string {
	var post string
	for k, v := range params {
		post += k + "=" + v + "&"
	}
	url = strings.TrimSuffix(url, "&")

	return h.Post(url, post)
}

// post发送字符串
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

// post提交表单
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

// post发送文件
func PostFile(url string, filePath string, params map[string]string) string {

	//打开文件句柄操作
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//创建一个模拟的form中的一个选项,这个form项现在是空的
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作, 设置文件的上传参数叫uploadfile, 文件名是filename,
	//相当于现在还没选择文件, form项里选择文件的选项
	fileWriter, err := bodyWriter.CreateFormFile("filename", path.Base(filePath))
	if err != nil {
		fmt.Println("error writing to buffer")
		panic(err)
	}

	//iocopy 这里相当于选择了文件,将文件放到form中
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		panic(err)
	}

	//获取上传文件的类型,multipart/form-data; boundary=...
	contentType := bodyWriter.FormDataContentType()

	//这个很关键,必须这样写关闭,不能使用defer关闭,不然会导致错误
	bodyWriter.Close()

	//这里就是上传的其他参数设置,可以使用 bodyWriter.WriteField(key, val) 方法
	//也可以自己在重新使用  multipart.NewWriter 重新建立一项,这个再server 会有例子
	if len(params) > 0 {
		//这种设置值得仿佛 和下面再从新创建一个的一样
		for key, val := range params {
			_ = bodyWriter.WriteField(key, val)
		}
	}

	//发送post请求到服务端
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(resp_body)
}

// 获取cookie里面的ticker
func (h *Http) GetTicket(uri string) (string, error) {
	uriss, _ := url.Parse(uri)
	cookies := h.Client.Jar.Cookies(uriss)
	for _, v := range cookies {
		if v.Name == "webwx_data_ticket" {
			return v.Value, nil
		}
	}
	return "", errors.New("not find 'webwx_data_ticket'")
}

// 上传资源文件
// 大于 1048576 使用秒传
// API_webwxpreview + "?fun=preview&mediaid=" + t.MediaId 预览
func (h *Http) UploadMedia(user *WxLoginStatus, username string, file string) string {

	uri := user.fileUri + "/webwxuploadmedia?f=json"
	// 检查文件是否存在
	if !utils.IsDirExist(file) {
		panic(errors.New("文件未找到"))
	}
	mediaCount++

	fileInfo, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	file_mime := mime.TypeByExtension(path.Ext(file))
	file_type := getFileType(file)
	file_size := fileInfo.Size()

	file_mod_time := fileInfo.ModTime()
	file_md5, _ := utils.Md5SumFile(file)
	ticket, _ := h.GetTicket(user.baseUri)

	uploadmediarequest := struct {
		BaseRequest   baseRequest
		ClientMediaId int
		TotalLen      int
		StartPos      int
		DataLen       int
		MediaType     int
		UploadType    int
		FromUserName  string
		ToUserName    string
		FileMd5       string
	}{
		BaseRequest:   user.BaseRequest,
		ClientMediaId: int(time.Now().Unix()),
		TotalLen:      int(file_size),
		DataLen:       int(file_size),
		StartPos:      0,
		MediaType:     4,
		UploadType:    2,
		FromUserName:  user.LoginUser.UserName,
		ToUserName:    username,
		FileMd5:       fmt.Sprintf("%x", file_md5),
	}

	bs, err := json.Marshal(uploadmediarequest)
	if err != nil {
		panic(err)
	}

	v := make(map[string]string)
	v["id"] = "WU_FILE_" + strconv.Itoa(mediaCount)
	v["name"] = path.Base(file)
	v["type"] = file_mime
	v["size"] = strconv.FormatInt(file_size, 10)
	v["mediatype"] = file_type
	v["uploadmediarequest"] = string(bs)
	v["lastModifieDate"] = file_mod_time.Format("Mon Jan 02 2006 15:04:15 GMT+0800 (CST)")
	v["webwx_data_ticket"] = ticket
	v["pass_ticket"] = user.passTicket
	v["filename"] = path.Base(file)

	content := PostFile(uri, file, v)

	ret := struct {
		BaseResponse      BaseResponse
		CDNThumbImgHeight int
		CDNThumbImgWidth  int
		MediaId           string
		StartPos          int
	}{}

	err = json.Unmarshal([]byte(content), &ret)
	if err != nil {
		// json解析错误
	}

	return ret.MediaId
}

// 获取上传文件类型
func getFileType(file string) string {
	switch path.Ext(file) {
	case ".jpg":
		return "pic"
	case ".png":
		return "pic"
	case ".gif":
		return "pic"
	case ".mp4":
		return "video"
	default:
		return "doc"
	}
}

// 下载图片，表情到本地
// 表情同样走这个接口
func (h *Http) DownMsgImg(user *WxLoginStatus, msgid string, file string) {
	uri := fmt.Sprintf(user.baseUri+"/webwxgetmsgimg?MsgID=%s&skey=%s", msgid, user.BaseRequest.Skey)
	res, err := h.Client.Get(uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s", uri, err)
		return
	}
	//TODO::保存
	if !utils.IsDirExist(file) {
		os.MkdirAll(file, 0755)
		fmt.Printf("dir %s created\n", file)
	}
	//根据URL文件名创建文件
	resp_body, err := ioutil.ReadAll(res.Body)
	ioutil.WriteFile(file, resp_body, os.ModePerm)
}

// 下载语音消息
func (h *Http) DownVoiceImg(user *WxLoginStatus, msgid string, file string) {
	uri := fmt.Sprintf(user.baseUri+"/webwxgetvoice?msgid=%s&skey=%s", msgid, user.BaseRequest.Skey)
	res, err := h.Client.Get(uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s", uri, err)
		return
	}
	//TODO::保存
	if !utils.IsDirExist(file) {
		os.MkdirAll(file, 0755)
		fmt.Printf("dir %s created\n", file)
	}
	//根据URL文件名创建文件
	resp_body, err := ioutil.ReadAll(res.Body)
	ioutil.WriteFile(file, resp_body, os.ModePerm)
}

// 下载语音消息
func (h *Http) DownVideoImg(user *WxLoginStatus, msgid string, file string) {
	uri := fmt.Sprintf(user.baseUri+"/webwxgetvideo?MsgID=%s&skey=%s", msgid, user.BaseRequest.Skey)

	res, err := h.Client.Get(uri)

	req, err := http.NewRequest("GET", uri, strings.NewReader(""))
	if err != nil {
		// handle error
	}
	req.Header.Set("Range", "bytes=0-")

	resp, err := h.Client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s", uri, err)
		return
	}
	//TODO::保存
	if !utils.IsDirExist(file) {
		os.MkdirAll(file, 0755)
		fmt.Printf("dir %s created\n", file)
	}
	//根据URL文件名创建文件
	resp_body, err := ioutil.ReadAll(res.Body)
	ioutil.WriteFile(file, resp_body, os.ModePerm)
}

// 下载文档保存到本地
func (h *Http) DownMsgFile(user *WxLoginStatus, msgid string, formUserName string, fileName string, file string) {

	baseUri := user.baseUri + "/webwxgetmedia?"
	ticket, _ := h.GetTicket(user.baseUri)
	uri := fmt.Sprintf("sender=%s&mediaid=%s&filename=%s&fromuser=%s&pass_ticket=%s&webwx_data_ticket=%s",
		formUserName, msgid, fileName, user.LoginUser.UserName, user.passTicket, ticket)
	uri = baseUri + uri

	res, err := h.Client.Get(uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s", uri, err)
		return
	}
	//根据URL文件名创建文件
	resp_body, err := ioutil.ReadAll(res.Body)
	ioutil.WriteFile(file, resp_body, os.ModePerm)
}

func (h *Http) SaveImage(url string, file string) {
	res, err := h.Client.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("%d HTTP ERROR:%s", url, err)
		return
	}
	//按分辨率目录保存图片
	Dirname := "tmp/"
	if !utils.IsDirExist(Dirname) {
		os.MkdirAll(Dirname, 0755)
		fmt.Printf("dir %s created\n", Dirname)
	}
	resp_body, err := ioutil.ReadAll(res.Body)
	//根据URL文件名创建文件
	ioutil.WriteFile(file, resp_body, os.ModePerm)
}
