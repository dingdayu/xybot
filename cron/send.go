package cron

import (
	"encoding/json"
	"fmt"

	"strconv"
	"time"
)

type Msg struct {
	Type         int
	EmojiFlag    int
	MediaId      string
	FromUserName string
	ToUserName   string
	LocalID      string
	ClientMsgId  string
	Content      string
}

type SendMsgContent struct {
	BaseRequest baseRequest
	Msg         Msg
	Scene       int
}

type RequestMsg struct {
	BaseResponse BaseResponse
	LocalID      string
	MsgID        string
}

// 发送文本消息
func (user *WxLoginStatus) SendTextMsg(username string, msg string) {
	url := fmt.Sprintf(user.baseUri+"/webwxsendmsg?pass_ticket=%s", user.passTicket)
	rand := rands(17)
	postData := SendMsgContent{
		BaseRequest: user.BaseRequest,
		Msg: Msg{
			Type:         1,
			Content:      msg,
			FromUserName: user.LoginUser.UserName,
			ToUserName:   username,
			LocalID:      rand,
			ClientMsgId:  rand,
		},
		Scene: 1,
	}

	bs, err := json.Marshal(postData)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(url, string(bs))
	var textReq RequestMsg
	err = json.Unmarshal([]byte(content), &textReq)
	if err != nil {
		// json解析错误
	}
}

// 发送图片消息
func (user *WxLoginStatus) SendImagesMsg(username string, file string) {

	uri := fmt.Sprintf(user.baseUri+"/webwxsendemoticon?fun=sys&f=json&pass_ticket=%s", user.passTicket)
	MediaId := NewHttp(user.uuid).UploadMedia(user, username, file)

	msg := SendMsgContent{
		BaseRequest: user.BaseRequest,
		Msg: Msg{
			Type:         47,
			EmojiFlag:    2,
			MediaId:      MediaId,
			FromUserName: user.LoginUser.UserName,
			ToUserName:   username,
			LocalID:      strconv.FormatInt(time.Now().Unix()*1e4, 10),
			ClientMsgId:  strconv.FormatInt(time.Now().Unix()*1e4, 10),
		},
		Scene: 0,
	}

	bs, err := json.Marshal(msg)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(uri, string(bs))
	var imgReq RequestMsg
	err = json.Unmarshal([]byte(content), &imgReq)
	if err != nil {
		// json解析错误
	}
}
