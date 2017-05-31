package cron

import (
	"encoding/json"
	"fmt"
	"github.com/IMQS/simplexml"
	"github.com/dingdayu/wxbot/model"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// 申请加为好友
	RequestFriend = "RequestFriend"
	//// 申请好友通过
	ReceiveAddFriendResult = "ReceiveAddFriendResult"
	//// 位置分享
	Location = "Location"
	//// 被踢下线
	//ReceiveKickOut  = "ReceiveKickOut"
	//// 登陆成功
	//LoginSucceed  = "LoginSucceed"
	//// 群成员信息更新
	//ReceiveMemberCardChanged  = "ReceiveMemberCardChanged"
	//// 新成员入群
	//NewMemberJoinCluster  = "NewMemberJoinCluster"




	// 修改群名称
	EditGroupName = "EditGroupName"
	// 新成员入群
	NewMemberJoinGroup  = "NewMemberJoinGroup"
	// 收到有人被踢出群
	MemberKickGroup  = "MemberKickGroup"
	// 收到有人退出群
	MemberExitGroup  = "MemberExitGroup"
	// 加入新群
	MeJoinGroup  = "MeJoinGroup"

	// 共享名片
	ShareCard = "ShareCard"
	// 分享的网页
	ShareWebPage = "ShareWebPage"
	// 分享小程序
	ShareApplet = "ShareApplet"
	// 转账
	Transfer = "Transfer"
	// 红包
	RedPackets = "RedPackets"
	// 撤回消息
	ReceiveMessageRevoke = "ReceiveMessageRevoke"

	ReceiveTextMessage  = "ReceiveTextMessage"
	ReceiveImageMessage = "ReceiveImageMessage"
	ReceiveVoiceMessage = "ReceiveVoiceMessage"
	ReceiveVideoMessage = "ReceiveVideoMessage"

	// 输入状态
	ReceiveInputState = "ReceiveInputState"
	// TODO:: ReceiveSignatureChanged 修改签名
	// TODO::  位置共享
	// TODO::  停止位置共享
)

func init() {

}

// 同步返回的消息体
type SyncStruct struct {
	BaseResponse           BaseResponse
	AddMsgCount            int
	AddMsgList             []types.Message
	ContinueFlag           int
	DelContactCount        int
	DelContactList         []model.Contact
	ModChatRoomMemberCount int
	ModChatRoomMemberList  []interface{}
	ModContactCount        int
	ModContactList         []model.Contact
	Profile                types.Profile
	SKey                   string
	SyncKey                SyncKey
	SyncCheckKey           SyncKey
}

func init() {
	go forcheck()
}

func forcheck() {
	for {
		for k := range WxMap {
			if item, ok := WxMap[k]; ok {
				if item.SyncOff {
					go item.CheckSync()
				}
			}

		}
		time.Sleep(3e9)

	}
}

func (user *WxLoginStatus) CheckSync() {
	//log.Println(user.uuid + ": 同步信息")
	// 锁
	user.SyncOff = false

	url := user.baseUri + "/synccheck?"

	data := make(map[string]string)
	data["r"] = strconv.FormatInt(time.Now().Unix(), 10)
	data["sid"] = user.BaseRequest.Sid
	data["uin"] = strconv.Itoa(user.BaseRequest.Uin)
	data["skey"] = user.BaseRequest.Skey
	data["deviceid"] = user.BaseRequest.DeviceID
	data["synckey"] = user.SyncKeyStr
	data["_"] = strconv.FormatInt(time.Now().Unix(), 10)

	content := NewHttp(user.uuid).Get(url, data)
	if len(content) <= 0 {
		user.SyncOff = true
		return
	}
	ret := utils.PregMatch(`window.synccheck=\{retcode:"(\d+)",selector:"(\d+)"\}`, content)

	retcode := ret[0]
	selector := ret[1]
	if retcode == "1100" || retcode == "1101" {
		log.Println("[" + user.uuid + "] 微信客户端正常退出")
		delete(WxMap, user.uuid)
		return
	}
	if retcode == "0" {
		user.handleCheckSync(selector)
		user.SyncOff = true
		return
	}

	if retcode == "1101" {
		fmt.Println("从其它设备上登了网页微信！")
		delete(WxMap, user.uuid)
		return
	}
	if retcode == "1100" {
		fmt.Println("从微信客户端上登出！")
		delete(WxMap, user.uuid)
		return
	}
}

func (user *WxLoginStatus) handleCheckSync(selector string) {
	if selector == "0" {
		return
	}
	// == 4 联系人修改资料
	// == 2 有新消息
	time.Sleep(2e9)
	user.Sync()

}

// 获取新消息
func (user *WxLoginStatus) Sync() {
	log.Println("[" + user.uuid + "] 有新消息")
	url := fmt.Sprintf(user.baseUri+"/webwxsync?sid=%s&skey=%s&lang=zh_CN&pass_ticket=%s", user.BaseRequest.Sid, user.BaseRequest.Skey, user.passTicket)

	type postDataStruct struct {
		BaseRequest baseRequest
		SyncKey     SyncKey
		Rr          string `json:"rr"`
	}
	var postData *postDataStruct = &postDataStruct{
		BaseRequest: user.BaseRequest,
		SyncKey:     user.SyncKey,
		Rr:          strconv.FormatInt(time.Now().Unix(), 10),
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, string(bs))
	var SyncMessage SyncStruct
	fmt.Println(content)
	err = json.Unmarshal([]byte(content), &SyncMessage)
	if err != nil {
		// json解析错误
		fmt.Println(err.Error())
	}

	// 更新SyncKey
	if SyncMessage.BaseResponse.Ret == 0 {
		user.SyncKey = SyncMessage.SyncKey
		user.SyncKeyStr = generateSyncKey(SyncMessage.SyncKey)

		user.handleSync(SyncMessage)
	}

}

// 处理新消息
func (user *WxLoginStatus) handleSync(SyncMessage SyncStruct) {
	if len(SyncMessage.ModContactList) > 0 {
		// 群变动
		// 检查 UserName 两个@@ 群成员变动
		// 否则 群成员编号
		for _, modContac := range SyncMessage.ModContactList {
			if strings.Contains(modContac.UserName, "@@") {
				if len(modContac.MemberList) > 0 {
					// 更新群成员信息, 组合usernmae ，然后请求群成员接口，然后更新群成员信息
				}
			}
			model.UpsertContact(&modContac)
		}
	}
	if len(SyncMessage.AddMsgList) > 0 {
		for _, msg := range SyncMessage.AddMsgList {
			user.handleMessage(msg)
		}
	}
	if len(SyncMessage.DelContactList) > 0 {
		for _, modContac := range SyncMessage.DelContactList {
			// 删除好友
			fmt.Println("删除好友：" + modContac.UserName)
		}
	}
}

// 处理消息类型
// AppMsgType app分享
// FromUserName 两个@@ 就是群消息
func (user *WxLoginStatus) handleMessage(Msg types.Message) {
	Msg.Content = FormatContent(Msg.Content)

	Contact := model.GetContactByUsername(Msg.FromUserName)

	log.Println("[" + Msg.ToUserName + "] 有来自[" + Msg.FromUserName + "]消息：[" + strconv.Itoa(Msg.MsgType) + "] " + Msg.Content)
	switch Msg.MsgType {
	case 1:

		if strings.Contains(Msg.Content, "webwxgetpubliclinkimg") && Msg.Url != "" {
			// 地理位置消息
			str := strings.Split(Msg.Content, ":\n")
			Msg.Content = str[0]
			//TODO::XML 在 `OriContent` 中
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        Location,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      ParseXml(Msg.OriContent, "location"),
				Url:          Msg.Url,
				SendTime:     Msg.CreateTime,
			}
			// 检查
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
			return
			//return msg
		}
		//TODO::如果FromUserName 存在于联系人，且
		if strings.Contains(Msg.Content, "过了你的朋友验证请求") && Contact.Id.Valid() {
			// 通过好友验证消息
			// 上面先处理的联系人变更， 所以，只要FromeUserName 能搜索到且搜到到字符就是新好友
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        ReceiveAddFriendResult,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      Msg.Content,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)
			log.Println(msg)
			return
		}

		// 文本消息
		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ReceiveTextMessage,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      Msg.Content,
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)
		log.Println(msg)
		//return msg
	case 3:

		// 图片消息
		img := ParseXml(Msg.Content, "img")
		md5 := img["md5"]

		// 下载文件
		// TODO::先检查文件是否已存在
		file := "./tmp/msg/img/" + md5 + ".jpg"
		if utils.IsDirExist(file) {
			NewHttp(user.uuid).DownImgMsg(user, Msg.MsgId, file)
		}

		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ReceiveImageMessage,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      img,
			Url:          "./tmp/msg/img/" + md5 + ".jpg",
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)

		log.Println(msg)
		return
	case 34:
		// 语音消息
		voicemsg := ParseXml(Msg.Content, "voicemsg")

		// 下载文件
		NewHttp(user.uuid).DownVoiceMsg(user, Msg.MsgId, Msg.MsgId+".mp3")

		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ReceiveVoiceMessage,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      voicemsg,
			Url:          "./tmp/msg/voice/" + Msg.MsgId + ".jpg",
			SendTime:     Msg.CreateTime,
		}
		log.Println(msg)
		return
	case 37:
		// 好友验证
		RequestMsg := ParseXml(Msg.Content, "msg")

		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        RequestFriend,
			FromUserName: Msg.RecommendInfo.UserName,
			FromNickName: Msg.RecommendInfo.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      RequestMsg,
			Url:          "./tmp/msg/voice/" + Msg.MsgId + ".jpg",
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)

		log.Println(msg)
		return
	case 42:
		// 共享名片
		RequestMsg := ParseXml(Msg.Content, "msg")

		RequestMsg["UserName"] = Msg.RecommendInfo.UserName
		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ShareCard,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      RequestMsg,
			Url:          "./tmp/msg/voice/" + Msg.MsgId + ".jpg",
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)

		log.Println(msg)
		return
	case 43:
		// 视频消息

		videomsg := ParseXml(Msg.Content, "videomsg")
		md5 := videomsg["md5"]

		// 下载文件
		// TODO::先检查文件是否已存在
		file := "./tmp/msg/video/" + md5 + ".jpg"
		if utils.IsDirExist(file) {
			NewHttp(user.uuid).DownImgMsg(user, Msg.MsgId, file)
		}

		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ReceiveVideoMessage,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      videomsg,
			Url:          file,
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)

		log.Println(msg)
		return
	case 47:
		// 动画表情
		// TODO:: 1、下载表情 msgid.gif  表情事件
		// 下载文件

		file := ""
		if Msg.Content == "" {
			file = "./tmp/msg/biaoqing/" + Msg.MsgId + ".jpg"
		} else {
			emoji := ParseXml(Msg.Content, "emoji")
			md5 := emoji["md5"]
			file = "./tmp/msg/emoji/" + md5 + ".jpg"
		}

		if utils.IsDirExist(file) {
			NewHttp(user.uuid).DownImgMsg(user, Msg.MsgId, file)
		}

		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ReceiveVideoMessage,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      "",
			Url:          file,
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)

		log.Println(msg)
	case 49:

		appmsg := ParseXml(Msg.Content, "appmsg")

		if Msg.Status == 3 && Msg.FileName == "微信转账" {
			// 微信转帐
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        Transfer,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      appmsg,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
			return
		}
		if Msg.Content == "该类型暂不支持，请在手机上查看" {
			return
		}

		if appmsg["type"] == "5" {
			// 分享的网页
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        ShareWebPage,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      appmsg,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
			return
		}

		if appmsg["type"] == "6" {
			// 分享的文件
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        ShareWebPage,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      appmsg,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			// todo::下载文件  md5 去重

			log.Println(msg)
			return
		}
		if appmsg["type"] == "33" {
			// 分享的小程序
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        ShareApplet,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      appmsg,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
			return
		}

		// 分享的网页 ,解析xml ， type 6 文件；8 动图； 33 小程序； 查询公众号FormUserName，如果在公众号，就是公众号，否则就是网页
		//<msg><appmsg appid="" sdkver="0"><type>17</type><title><![CDATA[我发起了位置共享]]></title></appmsg><fromusername>dingxiaoyu_ddy</fromusername></msg>

	case 51:
		if Msg.ToUserName == Msg.StatusNotifyUserName {
			// 点击事件（好友正在输入） 打开聊天框
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        ReceiveInputState,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      "",
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
			return
		}
	case 53:
		// 视频电话

	case 62:
		// 视频消息----

	case 10002:
		// 撤回消息
		sysmsg := ParseXml(Msg.Content, "sysmsg")
		msg := Message{
			MsgId:        Msg.MsgId,
			UUID:         user.uuid,
			Event:        ReceiveMessageRevoke,
			FromUserName: Msg.FromUserName,
			FromNickName: Contact.NickName,
			ToUserName:   Msg.ToUserName,
			Content:      sysmsg,
			SendTime:     Msg.CreateTime,
		}
		checkGroupMsg(&msg, user.LoginUser.UserName)

		log.Println(msg)
		return

	case 10000:
		if strings.Contains(Msg.Content, "利是") || strings.Contains(Msg.Content, "红包") {
			// 红包消息
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        RedPackets,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      Msg.Content,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
			return
		} else if strings.Contains(Msg.Content, "添加") && strings.Contains(Msg.Content, "打招呼") {
			// 通过好友申请，打招呼
			// todo::你已添加了小雨，现在可以开始聊天了。
		} else if strings.Contains(Msg.Content, "添加") && strings.Contains(Msg.Content, "聊天") {
			// 添加好友之间被通过，删除后重新添加
		} else if strings.Contains(Msg.Content, "加入了群聊") || strings.Contains(Msg.Content, "移出了群聊") || strings.Contains(Msg.Content, "改群名为") || strings.Contains(Msg.Content, "邀请你") || strings.Contains(Msg.Content, "分享的二维码加入群聊") {
			// 群成员改变，添加或移除
			if strings.Contains(Msg.Content, "邀请你") {
				// INVITE 邀请你加入新群
				group := model.GetContactByUsername(Msg.FromUserName)
				msg := Message{
					MsgId:        Msg.MsgId,
					UUID:         user.uuid,
					Event:        MeJoinGroup,
					FromUserName: "system",
					FromNickName: "system",
					GroupNickName: group.NickName,
					GroupUserName: group.UserName,
					ToUserName:   Msg.ToUserName,
					Content:      Msg.Content,
					SendTime:     Msg.CreateTime,
				}

				log.Println(msg)
				return
			}
			if strings.Contains(Msg.Content, "加入了群聊") || strings.Contains(Msg.Content, "分享的二维码加入群聊") {
				// ADD 新人入群
				name := utils.PregMatch(`邀请"(.+)"加入了群聊`, Msg.Content)
				if len(name) <= 0 {
					name = utils.PregMatch(`"(.+)"通过扫描.+分享的二维码加入群聊`, Msg.Content)
				}
				msg := Message{
					MsgId:        Msg.MsgId,
					UUID:         user.uuid,
					Event:        MeJoinGroup,
					FromUserName: Msg.FromUserName,
					FromNickName: Contact.NickName,
					ToUserName:   Msg.ToUserName,
					Content:      Msg.Content,
					SendTime:     Msg.CreateTime,
				}

				log.Println(msg)
				return


				fmt.Println(name)
			}
			if strings.Contains(Msg.Content, "移出了群聊") {
				// REMOVE 被移除
			}
			if strings.Contains(Msg.Content, "改群名为") {
				// RENAME 群名修改
				name := utils.PregMatch(`改群名为“(.+)”`, Msg.Content)
				fmt.Println(name)
			}
			if strings.Contains(Msg.Content, "移出群聊") {
				// BE_REMOVE 被移除群聊
			}
		}

		// 位置共享结束
	}
}

// 解析xml
func ParseXml(input string, name string) map[string]string {
	input, _ = FormatXml(input)
	rte, _ := simplexml.NewDocumentFromReader(strings.NewReader(input))
	att := rte.Root().Search().ByName(name).One().Attributes

	ma := map[string]string{}
	for _, v := range att {
		ma[v.Name] = v.Value
	}

	return ma
}

// 格式化xml
func FormatXml(input string) (string, string) {
	input = strings.Replace(input, "\\t", "", -1)
	input = strings.Replace(input, "\\", "", -1)

	input = FormatContent(input)

	username := ""
	content := utils.PregMatch(`(@\S+:)\n`, input)
	if len(content) > 0 {
		username = strings.Trim(content[0], ":")
		input = strings.Replace(input, content[0], "", 1)
	}

	input = strings.Replace(input, "\n", "", -1)

	return input, username
}

// 替换消息体
func FormatContent(content string) string {
	// 替换回车
	content = strings.Replace(content, "<br/>", "\n", -1)
	// 将html转义实例化
	content = html.UnescapeString(content)
	content = EmojiHandle(content)

	return content
}

// 处理 Emoji 表情
func EmojiHandle(content string) string {
	emoji := utils.PregMatch(`<span class="emoji emoji(.{1,10})"><\/span>`, content)
	if len(emoji) <= 0 {
		return content
	}
	re := regexp.MustCompile(`<span class="emoji emoji(.{1,10})"><\/span>`)
	src := re.FindAllString(content, -1)
	for k, v := range emoji {
		emjo := html.UnescapeString("&#x" + v + ";")
		content = strings.Replace(content, src[k], emjo, -1)
	}
	return content
}

type Message struct {
	MsgId         string `json:"msgid"`
	UUID          string
	Event         string
	FromUserName  string
	FromNickName  string
	ToUserName    string
	GroupUserName string
	GroupNickName string
	Content       interface{}
	Url           string `json:"url,omitempty"`
	SendTime      int
}

// 检查是否群消息，并补充群信息
func checkGroupMsg(msg *Message, toUserName string) {
	if strings.Contains(msg.ToUserName, "@@") {
		group := model.GetContactByUsername(msg.ToUserName)
		msg.GroupUserName = msg.ToUserName
		msg.GroupNickName = group.NickName
		msg.ToUserName = toUserName
	}
}
