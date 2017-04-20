package cron

import (
	"encoding/json"
	"fmt"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var host = ""

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

func CheckSync() string {

	url := host + "/synccheck?"

	data := make(map[string]string)
	data["r"] = strconv.FormatInt(time.Now().Unix(), 10)
	data["sid"] = ""
	data["uin"] = ""
	data["skey"] = ""
	data["deviceid"] = ""
	data["synckey"] = ""
	data["_"] = strconv.FormatInt(time.Now().Unix(), 10)

	content := NewHttp("").Get(url, data)

	ret := utils.PregMatch(`window.synccheck=\{retcode:"(\d+)",selector:"(\d+)"\}`, content)

	retcode := ret[1]
	selector := ret[2]
	if retcode == 1100 || retcode == 1101 {
		fmt.Println("微信客户端正常退出")
	}
	if retcode == 0 {
		handleCheckSync(selector)
	} else {
		fmt.Println("微信异常退出！")
	}

	return content
}

func handleCheckSync(selector int) {
	if selector == 0 {
		return
	}
	// == 4 联系人修改资料
}

func (user WxLoginStatus) Sync() {
	url := fmt.Sprintf(user.baseUri+"/webwxsync?sid=%s&skey=%s&lang=zh_CN&pass_ticket=%s", user.sid, user.skey, user.passTicket)

	type postDataStruct struct {
		BaseRequest baseRequest
		SyncKey     string
		rr          string
	}
	var postData *postDataStruct = &postDataStruct{
		BaseRequest: user.BaseRequest,
		SyncKey:     user.SyncKeyStr,
		rr:          string(time.Now().Unix()),
	}
	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, bs)
	var SyncMessage SyncStruct
	err = json.Unmarshal(byte(content), &SyncMessage)
	if err != nil {
		// json解析错误
		fmt.Println(err.Error())
	}
	handleSync(SyncMessage)
}

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

func handleMessage(Msg types.Message) {
	switch Msg.MsgType {
	case 1:

		if strings.Contains(Msg.Content, "webwxgetpubliclinkimg") && Msg.Url != "" {
			// 地理位置消息
		}

		// 通过好友验证消息
		// 文本消息

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
		// 微信转帐
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
