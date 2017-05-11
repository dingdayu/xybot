package cron

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"log"
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
func (user *WxLoginStatus) SendTextMsg(username string, msg string) (string, error) {
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
		log.Println("[" + user.uuid + "] [ERROR] [21001] generate json error")
		return "", errors.New("[21001] generate json error")
	}
	content := NewHttp(user.uuid).Post(url, string(bs))
	var textReq RequestMsg
	err = json.Unmarshal([]byte(content), &textReq)
	if err != nil {
		// json解析错误
		log.Println("[" + user.uuid + "] [ERROR] [21002] json error")
		return "", errors.New("[21002] json error")
	}
	log.Println("[" + user.uuid + "] 发送给：[" + username + "] :" + content)
	return textReq.MsgID, nil
}

// 发送图片消息
func (user *WxLoginStatus) SendEmoticonMsg(username string, file string) {

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

// 发送图片消息
func (user *WxLoginStatus) SendImagesMsg(username string, file string) {

	uri := fmt.Sprintf(user.baseUri+"/webwxsendmsgimg?fun=async&f=json&pass_ticket=%s", user.passTicket)
	MediaId := NewHttp(user.uuid).UploadMedia(user, username, file)

	msg := SendMsgContent{
		BaseRequest: user.BaseRequest,
		Msg: Msg{
			Type:         3,
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

// 发送视频消息
func (user *WxLoginStatus) SendVideoMsg(username string, file string) {

	uri := fmt.Sprintf(user.baseUri+"/webwxsendvideomsgx?fun=async&f=json&pass_ticket=%s", user.passTicket)
	MediaId := NewHttp(user.uuid).UploadMedia(user, username, file)

	msg := SendMsgContent{
		BaseRequest: user.BaseRequest,
		Msg: Msg{
			Type:         43,
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

// 通过验证消息
// 2 添加 3 通过
func (user *WxLoginStatus) VerifyUser(code string, msg string, username string, ticket string) bool {
	r := strconv.FormatInt(time.Now().Unix()*1e4, 13)
	uri := fmt.Sprintf(user.baseUri+"/webwxverifyuser?lang=zh_CN&r=%s&pass_ticket=%s", r, user.BaseRequest.Skey)

	type VierifyReq struct {
		BaseRequest        baseRequest
		Opcode             string
		VerifyUserListSize int
		VerifyUserList     map[string]string
		VerifyContent      string
		SceneListCount     int
		SceneList          []int
		skey               string
	}
	vierif := VierifyReq{
		user.BaseRequest,
		code,
		1,
		map[string]string{"Value": username, "VerifyUserTicket": ticket},
		msg,
		1,
		[]int{33},
		user.BaseRequest.Skey,
	}

	bs, err := json.Marshal(vierif)
	if err != nil {
		// json解析错误
	}
	content := NewHttp(user.uuid).Post(uri, string(bs))
	var req struct {
		BaseResponse BaseResponse
	}
	err = json.Unmarshal([]byte(content), &req)
	if err != nil {
		// json解析错误
	}
	if req.BaseResponse.Ret == 0 {
		return true
	} else {
		return false
	}
}
