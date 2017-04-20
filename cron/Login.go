package cron

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"os/user"
	"strconv"
	"time"
)

var redirectUri = ""
var fileUri = ""
var pushUri = ""
var baseUri = ""

// 登陆需要xml解析类型
type loginXml struct {
	Skey       string `xml:"skey"`
	Wxsid      string `xml:"wxsid"`
	Wxuin      string `xml:"wxuin"`
	PassTicket string `xml:"pass_ticket"`
}

type SyncKey struct {
	Count int
	List  []map[string]int
}

// 消息基础返回状态
type BaseResponse struct {
	ErrMsg string
	Ret    int
}

// 初始化登陆 返回消息体
type wxinitResponse struct {
	BaseResponse        BaseResponse
	User                types.User
	SyncKey             SyncKey
	SKey                string
	ClientVersion       int          // 客户端版本号 637863730
	ClickReportInterval int          // 单击间隔报告 600000
	ContactList         types.Member // 最近联系人
	Count               int          // 最近联系人个数
	InviteStartCount    int          // 翻译：邀请计数
	ChatSet             string
}

//
type WxLoginStatus struct {
	uuid        string
	skey        string
	sid         string
	uin         string
	passTicket  string
	BaseRequest baseRequest
	baseUri     string
	fileUri     string
	pushUri     string

	SyncKeyStr string
}

//wxinit 时需要提交的数据json格式
type baseRequest struct {
	Uin      int
	Sid      string
	Skey     string
	DeviceID string
}

// 待扫描队列
var uuidChannel chan string

func init() {
	uuidChannel = make(chan string)
	go check()
}

func Test() {
	done := make(chan struct{})

	fmt.Println("开始登陆:")
	uuid := GetUuid()
	generateQrCode(uuid)
	fmt.Println("开始等待扫描:")

	// 写入通道
	uuidChannel <- uuid

	<-done
	fmt.Println("登陆完成！")
}

func Xml() {
	xmlt := `<error><ret>0</ret><message></message><skey>@crypt_72e9aa7b_61a26569cadf952888c3f75936954c52</skey><wxsid>lpfwvW5zDt2cnskg</wxsid><wxuin>8104085</wxuin><pass_ticket>KA7GAP8T5UXXFnSj2YFlImPULlXwbFxQmGWvLtxLRBk%3D</pass_ticket><isgrayscale>1</isgrayscale></error>`
	var data loginXml
	err := xml.Unmarshal([]byte(xmlt), &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(data)

}

// 获取登陆uuid
func GetUuid() string {
	url := "https://login.weixin.qq.com/jslogin"
	data := make(map[string]string)
	data["appid"] = "wx782c26e4c19acffb"
	data["fun"] = "new"
	data["lang"] = "zh_CN"
	data["_"] = strconv.FormatInt(time.Now().Unix(), 10)

	content := NewHttp("").Get(url, data)
	fmt.Println(content)

	str := utils.PregMatch(`window.QRLogin.code = (\d+); window.QRLogin.uuid = \"(\S+?)\"`, content)
	if str == nil {
		fmt.Println("[ERROR] 请求UUID错误！")
		return ""
	}
	return str[2]
}

// 下载二维码
func generateQrCode(uuid string) {
	// 'https://login.weixin.qq.com/l/' . $this->uuid
	url := "https://login.weixin.qq.com/qrcode/" + uuid
	NewHttp("").SaveImage(url, uuid)
}

// 携程调用某账号的登陆检查
func check() {
	for {
		uuid := <-uuidChannel
		fmt.Println("开始检查状态:" + uuid)
		go waitForLogin(uuid)
	}
}

// 循环检查登陆
func waitForLogin(uuid string) {
	for i := 1; i < 10; i++ {
		tip := 1
		url := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%s&uuid=%s&_=%s", tip,
			uuid, strconv.FormatInt(time.Now().Unix(), 10))
		content := NewHttp(uuid).Get(url, make(map[string]string))

		code := utils.PregMatch(`window.code=(\d+);`, content)

		switch code[1] {
		case "201":
			fmt.Println("请点击微信登陆推送！")
			avater := utils.PregMatch(`(\S+?)window.userAvatar = '(\S+?)'`, content)
			fmt.Println(avater)
			tip = 0
		case "200":
			matches := utils.PregMatch(`window.redirect_uri="(https:\/\/(\S+?)\/\S+?)";`, content)

			redirectUri = matches[1] + "&fun=new"
			url := "https://%s/cgi-bin/mmwebwx-bin"
			fileUri = fmt.Sprintf(url, "file."+matches[2])
			pushUri = fmt.Sprintf(url, "webpush."+matches[2])
			baseUri = fmt.Sprintf(url, matches[2])
		case "408":
			fmt.Println("等待中……")
			tip = 1
			time.Sleep(time.Millisecond * 500)
		default:
			fmt.Println("登陆失败！错误码：" + code[0])
			tip = 1
		}
	}

}

// 获取登陆的第一次xml消息
func startLogin(uuid string, redirectUri string) {
	content := NewHttp(uuid).Get(redirectUri, make(map[string]string))
	var data loginXml
	err := xml.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(data)

}

// 获取个人信息及最近聊天
func (user WxLoginStatus) webwxinit(uuid string) {
	url := fmt.Sprintf(user.baseUri+"/webwxinit?r=%d", strconv.FormatInt(time.Now().Unix(), 10))

	//
	type postDataStruct struct {
		BaseRequest baseRequest
	}
	var postData *postDataStruct = &postDataStruct{
		BaseRequest: user.BaseRequest,
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(uuid).Post(url, bs)
	var wxLoginStatus WxLoginStatus
	err = json.Unmarshal(byte(content), &wxLoginStatus)
	if err != nil {
		// json解析错误
		fmt.Println(err.Error())
	}
	// TODO::将登陆状态信息放到map中
}

// 开启状态通知
func (user WxLoginStatus) statusNotify() string {
	url := fmt.Sprintf(user.baseUri+"/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", user.passTicket)
	type postDataStruct struct {
		BaseRequest  baseRequest
		Code         int
		FromUserName string
		ToUserName   string
		ClientMsgId  int
	}
	var postData *postDataStruct = &postDataStruct{
		BaseRequest:  user.BaseRequest,
		Code:         3,
		FromUserName: "", // 登陆用的username ，在 webwxinit中获得
		ToUserName:   "", // 同上
		ClientMsgId:  int(time.Now().Unix()),
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, bs)
	// TODO::默认不需要在处理信息了
	type statusNotifyRes struct {
		BaseResponse BaseResponse
		MsgID        string
	}
	var statusNotify statusNotifyRes
	err = json.Unmarshal(byte(content), &statusNotify)
	if err != nil {
		// json解析错误
	}
	return statusNotify.MsgID
}

// 拼接同步key
func generateSyncKey(synckey SyncKey) string {
	if len(synckey) > 0 {
		var syncString string
		for k, v := range synckey.List {
			synckey += k + "_" + v
		}
		return syncString
	}
	return ""
}

// 获取好友列表
func (user WxLoginStatus) getContactList(seq int) {
	url := fmt.Sprintf(user.baseUri+"/webwxgetcontact?pass_ticket=%s&skey=%s&r=%s&seq=%s", user.skey, user.passTicket,
		strconv.FormatInt(time.Now().Unix(), 10), string(seq))

	content := NewHttp(user.uuid).Post(url, "{}")
	// TODO::默认不需要在处理信息了
	type Members struct {
		BaseResponse BaseResponse
		MemberCount  int
		MemberList   []types.Member
		Seq          int
	}
	var members Members
	err := json.Unmarshal(byte(content), &members)
	if err != nil {
		// json解析错误
	}
	// TODO::处理好友列表
	if members.Seq != 0 {
		user.getContactList(members.Seq)
	}
}

// 获取群成员
func (user WxLoginStatus) getBatchGroupMembers() {
	url := fmt.Sprintf(user.baseUri+"/webwxbatchgetcontact?type=ex&r=%s&pass_ticket=%s", strconv.FormatInt(time.Now().Unix(), 10), user.passTicket)

	type postDataStruct struct {
		BaseRequest baseRequest
		Count       int
		List        []struct {
			EncryChatRoomId string
			UserName        string
		}
	}
	var postData *postDataStruct = &postDataStruct{
		BaseRequest: user.BaseRequest,
		Count:       0,
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, bs)
	// TODO::默认不需要在处理信息了
	type GroupMembers struct {
		BaseResponse BaseResponse
		ContactList  []types.Member
		Count        int
	}
	var groupMembers GroupMembers
	err = json.Unmarshal(byte(content), &groupMembers)
	if err != nil {
		// json解析错误
	}

}

//
//func webwxinit()  {
//	url := baseUrl + "/webwxinit?r=" + strconv.FormatInt(time.Now().Unix(), 10);
//
//}
//
//func statusNotify() {
//	baseUrl + "/webwxstatusnotify?lang=zh_CN&pass_ticket=" + passTicket;
//}
