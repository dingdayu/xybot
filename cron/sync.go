package cron

import (
	"encoding/json"
	"fmt"
	"github.com/IMQS/simplexml"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	DelContactList         []interface{}
	ModChatRoomMemberCount int
	ModChatRoomMemberList  []interface{}
	ModContactCount        int
	ModContactList         []types.ModContact
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
		for k, _ := range WxMap {
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
	log.Println(user.uuid + ": 同步信息")
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
	ret := utils.PregMatch(`window.synccheck=\{retcode:"(\d+)",selector:"(\d+)"\}`, content)

	retcode := ret[0]
	selector := ret[1]
	if retcode == "1100" || retcode == "1101" {
		fmt.Println("微信客户端正常退出")
		delete(WxMap, user.uuid)
	}
	if retcode == "0" {
		user.handleCheckSync(selector)
		user.SyncOff = true
	}

	if retcode == "1101" {
		fmt.Println("从其它设备上登了网页微信！")
		delete(WxMap, user.uuid)
	}
	if retcode == "1100" {
		fmt.Println("从微信客户端上登出！")
		delete(WxMap, user.uuid)
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
	err = json.Unmarshal([]byte(content), &SyncMessage)
	if err != nil {
		// json解析错误
		fmt.Println(err.Error())
	}

	// 更新SyncKey
	if SyncMessage.BaseResponse.Ret == 0 {
		user.SyncKey = SyncMessage.SyncKey
		user.SyncKeyStr = generateSyncKey(SyncMessage.SyncKey)

		handleSync(SyncMessage)
	}

}

// 处理新消息
func handleSync(SyncMessage SyncStruct) {
	if len(SyncMessage.ModContactList) > 0 {
		// 群变动
		// 检查 UserName 两个@@ 群成员变动
		// 否则 群成员编号
		for _, modContac := range SyncMessage.ModContactList {
			if strings.Contains(modContac.UserName, "@@") {
				// 更新群成员信息
				// 新建群

			} else {
				// 联系人更新资料 TODO::暂时无效
			}
		}

	}
	if len(SyncMessage.AddMsgList) > 0 {
		for _, msg := range SyncMessage.AddMsgList {
			handleMessage(msg)
		}
	}
}

// 处理消息类型
// AppMsgType app分享
// FromUserName 两个@@ 就是群消息
func handleMessage(Msg types.Message) {
	Msg.Content = FormatContent(Msg.Content)

	log.Println("[" + Msg.ToUserName + "] 有来自[" + Msg.FromUserName + "]消息：[" + strconv.Itoa(Msg.MsgType) + "] " + Msg.Content)
	switch Msg.MsgType {
	case 1:

		if strings.Contains(Msg.Content, "webwxgetpubliclinkimg") && Msg.Url != "" {
			// 地理位置消息
			str := strings.Split(Msg.Content, ":\n")
			Msg.Content = str[0]
			fmt.Println(Msg.Content)
			fmt.Println(Msg.Url)
		}
		//TODO::如果FromUserName 存在于联系人，且
		if strings.Contains(Msg.Content, "过了你的朋友验证请求") && Msg.FromUserName == "" {
			// 通过好友验证消息
			// 上面先处理的联系人变更， 所以，只要FromeUserName 能搜索到且搜到到字符就是新好友
		}

		// 文本消息
		fmt.Println(FormatContent(Msg.Content))

	case 3:
		// 图片消息

	case 34:
		// 语音消息

	case 37:
		// 好友验证
		// 提取好友头像
		matches := utils.PregMatch(`bigheadimgurl="(.+?)"`, Msg.Content)
		avatar := matches[1]
		Msg.RecommendInfo.NickName = FormatContent(Msg.RecommendInfo.NickName)
		fmt.Println(avatar)
		fmt.Println(Msg.RecommendInfo)
	case 42:
		// 共享名片
		fmt.Println(Msg)

	case 43:
		// 视频消息
	case 47:
		// 动画表情
		// TODO:: 1、下载表情 msgid.gif
	case 49:

		if Msg.Status == 3 && Msg.FileName == "微信转账" {
			// 微信转帐
		}
		if Msg.Content == "该类型暂不支持，请在手机上查看" {
			return
		}
		// 分享的网页 ,解析xml ， type 6 文件； 33 小程序； 查询公众号FormUserName，如果在公众号，就是公众号，否则就是网页

	case 51:
		if Msg.ToUserName == Msg.StatusNotifyUserName {
			// 点击事件（好友正在输入）
		}
	case 53:
		// 视频电话

	case 62:
		// 视频消息----

	case 10002:
		// 撤回消息
		msgID := utils.PregMatch(`<msgid>(\d+)<\/msgid>`, Msg.Content)
		fmt.Println(msgID)
		// 通过msgID获取上一条消息

	case 10000:
		if strings.Contains(Msg.Content, "利是") || strings.Contains(Msg.Content, "红包") {
			// 红包消息
		}
		if strings.Contains(Msg.Content, "添加") || strings.Contains(Msg.Content, "打招呼") {
			// 好友申请，打招呼
		}
		// 群成员改变，添加或移除
		GroupChange(Msg)

	}
}

/**

@da0e51d64fca0a686c08ee57a6d87e364c3a7506da8a7c2c128c4dd421e23f7d:
<br><msg><appmsg appid="" sdkver=""><title>å¾·å›½åŒåˆ©å•åˆ€ç™»é™†ä¸­å›½ï¼Œå°†åŽŸç”¨äºŽä¸­å›½åŒºçš„1äº¿å®£ä¼ å¹¿å‘Šè´¹ç›´æŽ¥åšæˆå®¢æˆ·å…è´¹ä½“éªŒã€‚</title><des>http://test.ebeq.net/index.php?m=home&amp;c=goods&amp;a=index&amp;gid=3&amp;uid=10303</des><action>view</action><type>5</type><showtype>0</showtype><content></content><url>http://test.ebeq.net/index.php?m=home&amp;c=goods&amp;a=index&amp;gid=3&amp;uid=10303</url><dataurl></dataurl><lowurl></lowurl><lowdataurl></lowdataurl><recorditem><![CDATA[]]></recorditem><thumburl>http://wx.qlogo.cn/mmopen/2jxblmcQQWyRfGrDib2G5ePlf8KVfj4R4ChTKcfHbaj9WWS2tsJqJpSQKGhpTtySMFrPCiaIrROMQrHD2LrqyFt4pmPRLiaQ07a/0</thumburl><extinfo></extinfo><sourceusername></sourceusername><sourcedisplayname></sourcedisplayname><commenturl></commenturl><appattach><totallen>0</totallen><attachid></attachid><emoticonmd5></emoticonmd5><fileext></fileext></appattach></appmsg><fromusername>wxid_ss7bwx1wixe722</fromusername><scene>0</scene><appinfo><version>1</version><appname></appname></appinfo><commenturl></commenturl></msg>
<br>

@da0e51d64fca0a686c08ee57a6d87e364c3a7506da8a7c2c128c4dd421e23f7d:
<br>wxid_ss7bwx1wixe722:
<br><?xml version="1.0"?>
<br><msg>
<br>    <videomsg aeskey="39376464363533356333376431333537" cdnthumbaeskey="39376464363533356333376431333537" cdnvideourl="304c0201000445304302010002042a31c66f02033d14b9020497e503b7020458f7f2440421777869645f656d61616463787a72726b7132323330355f313439323634343431380201000201000400" cdnthumburl="304c0201000445304302010002042a31c66f02033d14b9020497e503b7020458f7f2440421777869645f656d61616463787a72726b7132323330355f313439323634343431380201000201000400" length="7030184" playlength="103" cdnthumblength="6287" cdnthumbwidth="0" cdnthumbheight="0" fromusername="wxid_ss7bwx1wixe722" md5="0bbabb5d333fee84695d8087b1fe0553" newmd5="" isad="0" />
<br></msg>
<br>

*/

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
