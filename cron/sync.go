package cron

import (
	"encoding/json"
	"fmt"
	"github.com/dingdayu/wxbot/simplexml"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"html"
	"os"
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

func (user *WxLoginStatus) CheckSync() {
	fmt.Println("检查消息 =>" + time.Now().Format("2006-01-02 15:04"))
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
	}
	if retcode == "0" {
		user.handleCheckSync(selector)
	}

	if retcode == "1101" {
		fmt.Println(content)
		fmt.Println("从其它设备上登了网页微信！")
		os.Exit(0)
	}
	if retcode == "1100" {
		fmt.Println("从微信客户端上登出！")
		os.Exit(0)
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
	fmt.Println("新有消息 =>" + time.Now().Format("2006-01-02 15:04"))
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
			} else {
				// 联系人更新资料
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
	Msg.Content = strings.Replace(Msg.Content, "<br>", "\n", 99)
	switch Msg.MsgType {
	case 1:

		if strings.Contains(Msg.Content, "webwxgetpubliclinkimg") && Msg.Url != "" {
			// 地理位置消息
		}

		// 通过好友验证消息
		// 文本消息
		fmt.Println(FormatContent(Msg.Content))

	case 3:
		// 图片消息
	case 34:
		// 语音消息
	case 37:
		// 好友验证
	case 42:
		// 共享名片

	case 43:
		// 视频消息
	case 47:
		// 动画表情
	case 49:

		if Msg.Status == 3 && Msg.FileName == "微信转账" {
			// 微信转帐
		}
		if Msg.Content == "该类型暂不支持，请在手机上查看" {
			return
		}
		// 分享的网页
	case 51:
		// 点击事件（好友正在输入）
	case 53:
		// 视频电话

	case 62:
		// 视频消息----

	case 10002:
		// 撤回消息
	case 10000:
		// 红包消息
		// 好友申请，打招呼
		// 群成员改变，添加或移除

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

func parseXml(xml string) {
	if strings.HasSuffix(xml, "@") {
		content := utils.PregMatch(`(@\S+:\\n)`, xml)
		// create a document from a reader
		doc, err := simplexml.NewDocumentFromReader(strings.NewReader(content[1]))
		if err != nil {
			panic(err)
		}

		// get the fizz tag and value
		fizz := doc.Root().Search().ByName("foo").ByName("fizz").One()
		if fizz == nil {
			panic("fizz is missing")
		}

		fv, err := fizz.Value()
		if err != nil {
			panic(err)
		}

		fmt.Println("fizz: ", fv)
	}
}

func FormatContent(content string) string {
	// 替换回车
	content = strings.Replace(content, "<br/>", "\n", -1)
	// 将html转义实例化
	content = html.UnescapeString(content)

	emoji := utils.PregMatch(`<span class="emoji emoji(.{1,10})"><\/span>`, content)
	re := regexp.MustCompile(`<span class="emoji emoji(.{1,10})"><\/span>`)
	src := re.FindAllString(content, -1)

	for k, v := range emoji {
		emjo := html.UnescapeString("&#x" + v + ";")
		content = strings.Replace(content, src[k], emjo, -1)
	}
	return content
}
