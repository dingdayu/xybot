package types

type AppInfo struct {
	AppID string
	Type  int
}

// 推荐信息 [名片推荐]
type RecommendInfo struct {
	Alias      string
	AttrStatus int
	City       string
	Content    string
	NickName   string
	OpCode     int
	Province   string
	QQNum      int
	Scene      int
	Sex        int
	Signature  string
	Ticket     string
	UserName   string
	VerifyFlag int
}

// 同步信息时返回的Message
type Message struct {
	AppInfo              AppInfo
	AppMsgType           int
	Content              string
	CreateTime           int
	FileName             string
	FileSize             string
	ForwardFlag          int // 转发标记
	FromUserName         string
	HasProductId         int
	ImgHeight            int
	ImgStatus            int
	ImgWidth             int
	MediaId              string
	MsgId                string
	MsgType              int
	NewMsgId             int
	OriContent           string
	PlayLength           int
	RecommendInfo        RecommendInfo
	Status               int
	StatusNotifyCode     int
	StatusNotifyUserName string
	SubMsgType           int
	Ticket               string
	ToUserName           string
	Url                  string
	VoiceLength          int

	GroupUserName string
	GroupNickName string
}

type BatchGetContact struct {
	UserName        string
	EncryChatRoomId string
}

type SendMsg struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      string
	ClientMsgId  string
}
