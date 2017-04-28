package cron

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dingdayu/wxbot/model"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"strings"
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
	Wxuin      int    `xml:"wxuin"`
	PassTicket string `xml:"pass_ticket"`
}

type SyncKey struct {
	Count int
	List  []struct {
		Key int
		Val int
	}
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
	ClientVersion       int            // 客户端版本号 637863730
	ClickReportInterval int            // 单击间隔报告 600000
	ContactList         []types.Member // 最近联系人
	Count               int            // 最近联系人个数
	InviteStartCount    int            // 翻译：邀请计数
	ChatSet             string
	SystemTime          int
	GrayScale           int
	MPSubscribeMsgCount int
	MPSubscribeMsgList  []types.MPSubscribeMsg
}

//
type WxLoginStatus struct {
	uuid        string
	passTicket  string
	BaseRequest baseRequest
	LoginUser   types.User

	// 联系人
	ContactList []types.Member

	baseUri string
	fileUri string
	pushUri string

	SyncKey    SyncKey
	SyncKeyStr string
}

//wxinit 时需要提交的数据json格式
type baseRequest struct {
	Uin      int
	Sid      string
	Skey     string
	DeviceID string
}

// 获取详细信息接口返回的数据结构
type GroupMembers struct {
	BaseResponse BaseResponse
	ContactList  []types.Member
	Count        int
}

// 待扫描队列
var uuidChannel chan string

// 登陆的用户组
var WxMap map[string]*WxLoginStatus

func init() {
	uuidChannel = make(chan string)
	WxMap = make(map[string]*WxLoginStatus)
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

func Js() {
	js := `{
"BaseResponse": {
"Ret": 0,
"ErrMsg": ""
}
,
"Count": 2,
"ContactList": [{
"Uin": 0,
"UserName": "filehelper",
"NickName": "文件传输助手",
"HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?seq=0&username=filehelper&skey=@crypt_3cc836df_d87639179f41792d916c08038b2ffc50",
"ContactFlag": 0,
"MemberCount": 0,
"MemberList": [],
"RemarkName": "",
"HideInputBarFlag": 0,
"Sex": 0,
"Signature": "",
"VerifyFlag": 0,
"OwnerUin": 0,
"PYInitial": "WJCSZS",
"PYQuanPin": "wenjianchuanshuzhushou",
"RemarkPYInitial": "",
"RemarkPYQuanPin": "",
"StarFriend": 0,
"AppAccountFlag": 0,
"Statues": 0,
"AttrStatus": 0,
"Province": "",
"City": "",
"Alias": "",
"SnsFlag": 0,
"UniFriend": 0,
"DisplayName": "",
"ChatRoomId": 0,
"KeyWord": "fil",
"EncryChatRoomId": "",
"IsOwner": 0
}
,{
"Uin": 0,
"UserName": "weixin",
"NickName": "微信团队",
"HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?seq=658070034&username=weixin&skey=@crypt_3cc836df_d87639179f41792d916c08038b2ffc50",
"ContactFlag": 3,
"MemberCount": 0,
"MemberList": [],
"RemarkName": "",
"HideInputBarFlag": 0,
"Sex": 0,
"Signature": "微信团队官方帐号",
"VerifyFlag": 56,
"OwnerUin": 0,
"PYInitial": "WXTD",
"PYQuanPin": "weixintuandui",
"RemarkPYInitial": "",
"RemarkPYQuanPin": "",
"StarFriend": 0,
"AppAccountFlag": 0,
"Statues": 0,
"AttrStatus": 4,
"Province": "",
"City": "",
"Alias": "",
"SnsFlag": 0,
"UniFriend": 0,
"DisplayName": "",
"ChatRoomId": 0,
"KeyWord": "wei",
"EncryChatRoomId": "",
"IsOwner": 0
}
],
"SyncKey": {
"Count": 4,
"List": [{
"Key": 1,
"Val": 658070058
}
,{
"Key": 2,
"Val": 658070059
}
,{
"Key": 3,
"Val": 658070050
}
,{
"Key": 1000,
"Val": 1492994161
}
]
}
,
"User": {
"Uin": 2363862471,
"UserName": "@3b53de9219d5ed41affdbc0018fb7c529b1187ebcbb4de5ebf1bdc7ac99d67f3",
"NickName": "小雨6",
"HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?seq=1434526769&username=@3b53de9219d5ed41affdbc0018fb7c529b1187ebcbb4de5ebf1bdc7ac99d67f3&skey=@crypt_3cc836df_d87639179f41792d916c08038b2ffc50",
"RemarkName": "",
"PYInitial": "",
"PYQuanPin": "",
"RemarkPYInitial": "",
"RemarkPYQuanPin": "",
"HideInputBarFlag": 0,
"StarFriend": 0,
"Sex": 0,
"Signature": "",
"AppAccountFlag": 0,
"VerifyFlag": 0,
"ContactFlag": 0,
"WebWxPluginSwitch": 0,
"HeadImgFlag": 0,
"SnsFlag": 0
}
,
"ChatSet": "filehelper,weixin,",
"SKey": "@crypt_3cc836df_d87639179f41792d916c08038b2ffc50",
"ClientVersion": 637863730,
"SystemTime": 1493005587,
"GrayScale": 1,
"InviteStartCount": 40,
"MPSubscribeMsgCount": 0,
"MPSubscribeMsgList": [],
"ClickReportInterval": 600000
}`
	var wxinitResponse wxinitResponse
	err := json.Unmarshal([]byte(js), &wxinitResponse)
	if err != nil {
		// json解析错误
		fmt.Println(err.Error())
	}
	for _, item := range wxinitResponse.ContactList {
		fmt.Println(item)
		var contact = model.Contact{}
		utils.Struct2Struct(item, &contact)
		contact.LoginUin = wxinitResponse.User.Uin
		model.AddContact(contact)
	}
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
		tip := 0
		url := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%s&uuid=%s&_=%s", strconv.Itoa(tip),
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

			fmt.Println("开始请求xml")
			loginXm := startLogin(uuid, redirectUri)

			fmt.Println("开始拼接登陆状态")

			baseRequest := baseRequest{
				Uin:      loginXm.Wxuin,
				Sid:      loginXm.Wxsid,
				Skey:     loginXm.Skey,
				DeviceID: "e" + rands(15),
			}
			if v, ok := WxMap[uuid]; ok {
				v.uuid = uuid
				v.baseUri = baseUri
				v.pushUri = pushUri
				v.fileUri = fileUri
				v.passTicket = strings.TrimSpace(loginXm.PassTicket)
				v.BaseRequest = baseRequest
			} else {
				WxMap[uuid] = &WxLoginStatus{
					uuid:        uuid,
					baseUri:     baseUri,
					pushUri:     pushUri,
					fileUri:     fileUri,
					passTicket:  strings.TrimSpace(loginXm.PassTicket),
					BaseRequest: baseRequest,
				}
			}

			//fmt.Println(WxMap[uuid])
			WxMap[uuid].webwxinit(uuid)
			return

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
func startLogin(uuid string, redirectUri string) loginXml {
	content := NewHttp(uuid).Get(redirectUri, make(map[string]string))
	var data loginXml
	err := xml.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	return data

}

// 获取个人信息及最近聊天
func (user *WxLoginStatus) webwxinit(uuid string) {
	url := fmt.Sprintf(user.baseUri+"/webwxinit?r=%d&lang=zh_CN&pass_ticket=%s", strconv.FormatInt(time.Now().Unix(), 10), user.passTicket)

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

	content := NewHttp(uuid).Post(url, string(bs))
	//fmt.Println(content)
	var wxinitResponse wxinitResponse
	err = json.Unmarshal([]byte(content), &wxinitResponse)
	if err != nil {
		// json解析错误
		fmt.Println(err.Error())
	}

	if wxinitResponse.BaseResponse.Ret != 0 {
		fmt.Println("初始化失败！")
		return
	}

	// 将登陆状态信息放到map中
	user.SyncKey = wxinitResponse.SyncKey
	user.SyncKeyStr = generateSyncKey(wxinitResponse.SyncKey)

	// 保存登陆人资料
	fmt.Println("初始化个人资料")
	// 复制登陆资料
	user.LoginUser = wxinitResponse.User
	var dbUser = model.User{}
	utils.Struct2Struct(wxinitResponse.User, &dbUser)
	dbUser.UUID = uuid
	dbUser.Time = int(time.Now().Unix())
	model.UpsertUser(dbUser)

	// 初始化联系人
	fmt.Println("初始化联系人")
	// TODO::保存最近联系人里的群，公众号，联系人
	for _, item := range wxinitResponse.ContactList {
		var contact = model.Contact{}
		utils.Struct2Struct(item, &contact)
		contact.LoginUin = wxinitResponse.User.Uin
		contact.UUID = uuid
		contact.ContactType = getContactType(item, user.LoginUser.UserName)
		contact.HeadImgUrl = user.baseUri + item.HeadImgUrl
		model.UpsertContact(&contact)
	}
	// TODO::复制联系人对象
	// 获取全部的好友列表
	user.getContactList(0)

	// 根据群UserName获取所有群的群成员列表
	chatRommList := model.GetLimitContact(bson.M{"uuid": user.uuid, "contact_type": "ChatRooms"}, 50, 0)
	batch := []types.BatchGetContact{}
	for _, v := range *chatRommList {
		batch = append(batch, types.BatchGetContact{UserName: v.UserName, EncryChatRoomId: ""})
	}
	groupMembers := user.getBatchGroupMembers(batch)
	for _, item := range groupMembers.ContactList {
		var contact = model.Contact{}
		utils.Struct2Struct(item, &contact)
		contact.LoginUin = user.BaseRequest.Uin
		contact.UUID = user.uuid
		contact.HeadImgUrl = user.baseUri + item.HeadImgUrl
		contact.ContactType = getContactType(item, user.LoginUser.UserName)
		model.UpsertContact(&contact)
		// 将群成员加入成员表
		if len(item.MemberList) > 0 {
			for _, items := range item.MemberList {
				var member = model.Member{}
				utils.Struct2Struct(items, &member)
				member.HeadImgUrl = user.baseUri + member.HeadImgUrl
				member.ChatRoomUserName = item.UserName
				member.UUID = user.uuid
				member.LoginUin = user.BaseRequest.Uin
				model.UpsertMember(&member)
			}
		}
	}

	// TODO::获取所有群成员的资料信息
	fmt.Println("第二次 更新群成员资料")
	count := model.GetMemberCount(bson.M{"uuid": user.uuid})
	page := count / 50
	for i := 1; i < page; i++ {
		time.Sleep(2e9)
		fmt.Println("更新群成员资料：第" + strconv.Itoa(i) + "页,共" + strconv.Itoa(page) + "页")
		members := model.GetLimitMember(50, i*50)
		batch := []types.BatchGetContact{}
		for _, v := range *members {
			batch = append(batch, types.BatchGetContact{UserName: v.UserName, EncryChatRoomId: v.ChatRoomUserName})
		}
		groupMembers := user.getBatchGroupMembers(batch)
		for _, item := range groupMembers.ContactList {
			var contact = model.Member{}
			utils.Struct2Struct(item, &contact)
			contact.LoginUin = user.BaseRequest.Uin
			contact.UUID = user.uuid
			contact.HeadImgUrl = user.baseUri + item.HeadImgUrl
			model.UpsertMember(&contact)
		}
	}

	user.CheckSync()

}

// 开启状态通知
func (user *WxLoginStatus) statusNotify() string {
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
	content := NewHttp(user.uuid).Post(url, string(bs))
	// TODO::默认不需要在处理信息了
	type statusNotifyRes struct {
		BaseResponse BaseResponse
		MsgID        string
	}
	var statusNotify statusNotifyRes
	err = json.Unmarshal([]byte(content), &statusNotify)
	if err != nil {
		// json解析错误
	}
	return statusNotify.MsgID
}

// 拼接同步key
func generateSyncKey(synckey SyncKey) string {
	if len(synckey.List) > 0 {
		var syncString bytes.Buffer
		for _, v := range synckey.List {
			syncString.WriteString(string(v.Key) + "_" + string(v.Val))
		}
		return syncString.String()
	}
	return ""
}

// 获取好友列表
func (user *WxLoginStatus) getContactList(seq int) {
	fmt.Println("拉取好友列表")
	url := fmt.Sprintf(user.baseUri+"/webwxgetcontact?lang=zh_CN&pass_ticket=%s&r=%s&seq=%s&skey=%s", user.passTicket,
		strconv.FormatInt(time.Now().Unix(), 10), strconv.Itoa(seq), user.BaseRequest.Skey)

	content := NewHttp(user.uuid).Get(url, make(map[string]string))
	//
	type Members struct {
		BaseResponse BaseResponse
		MemberCount  int
		MemberList   []types.Member
		Seq          int
	}
	var members Members
	err := json.Unmarshal([]byte(content), &members)
	if err != nil {
		// json解析错误
	}

	// 初始化联系人
	for _, item := range members.MemberList {
		var contact = model.Contact{}
		utils.Struct2Struct(item, &contact)
		contact.LoginUin = user.BaseRequest.Uin
		contact.UUID = user.uuid
		contact.HeadImgUrl = user.baseUri + item.HeadImgUrl
		contact.ContactType = getContactType(item, user.LoginUser.UserName)
		model.UpsertContact(&contact)
	}
	fmt.Println(members.MemberCount)

	if members.Seq != 0 {
		user.getContactList(members.Seq)
	}
}

// 获取详细资料 （可用于公众号，好友，群等）
// 通常用于群资料和群成员资料
// 每次最多50个
func (user *WxLoginStatus) getBatchGroupMembers(batch []types.BatchGetContact) GroupMembers {

	if len(batch) > 50 {
		batch = batch[:49]
	}

	url := fmt.Sprintf(user.baseUri+"/webwxbatchgetcontact?type=ex&r=%s&pass_ticket=%s", strconv.FormatInt(time.Now().Unix(), 10), user.passTicket)

	type postDataStruct struct {
		BaseRequest baseRequest
		Count       int
		List        []types.BatchGetContact
	}
	var postData *postDataStruct = &postDataStruct{
		BaseRequest: user.BaseRequest,
		Count:       len(batch),
		List:        batch,
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, string(bs))
	//
	var groupMembers GroupMembers
	err = json.Unmarshal([]byte(content), &groupMembers)
	if err != nil {
		// json解析错误
	}
	return groupMembers
}

func (user *WxLoginStatus) accountLogout() {
	url := fmt.Sprintf(user.baseUri+"/webwxlogout?redirect=1&type=1&skey=%s", user.BaseRequest.Skey)

	type postDataStruct struct {
		Sid string
		Uin int
	}

	var postData *postDataStruct = &postDataStruct{
		Sid: user.BaseRequest.Sid,
		Uin: user.BaseRequest.Uin,
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, string(bs))
	//
	fmt.Println(content)
}

// 随机数字符串
func rands(n int) string {
	rand.Seed(int64(time.Now().Nanosecond()))

	str := strconv.Itoa(rand.Int())
	stra := []rune(str)
	return string(stra[:n])
}

// 获取账号类型
func getContactType(member types.Member, selfUserName string) string {
	// 检查系统号
	for _, v := range types.SPECIAL_USERS {
		if member.UserName == v {
			return "System"
		}
	}
	for _, v := range types.SPECIAL_USERS_NAME {
		if member.NickName == v {
			return "System"
		}
	}

	// 会有一些被注销的公众号，相比通讯录多出几个,企业号也在里面
	if member.VerifyFlag&8 != 0 {
		return "MP"
	}

	// 未保存到通讯录，但最近有过消息的依然会在这里面，相比通讯录会多出几个
	if strings.HasPrefix(member.UserName, "@@") {
		return "ChatRooms"
	}

	// 相比于通讯录 总好友数 会少一个，因为那是自己
	if selfUserName == member.UserName {
		return "Self"
	}

	return "Friends"
}
