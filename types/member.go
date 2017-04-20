package types

// 好友信息类型
type Member struct {
	Alias           string
	ContactFlag     int
	DisplayName     string
	ChatRoomId      int // 是否聊天室 0 否
	EncryChatRoomId string
	MemberCount     int        // 聊天室成员数量
	MemberList      RoomMember // 聊天室成员
	HeadImgUrl      string
	IsOwner         int // 是否所有者
	OwnerUin        int // 所有者id
	KeyWord         string

	NickName        string // 昵称
	PYInitial       string // 昵称简拼
	PYQuanPin       string // 昵称全拼
	Province        string // 省份
	City            string // 城市
	RemarkName      string // 备注
	RemarkPYInitial string
	RemarkPYQuanPin string
	Sex             int
	Signature       string
	SnsFlag         int
	StarFriend      int
	Statues         int
	Uin             int
	UniFriend       int
	UserName        string
	VerifyFlag      int // VerifyFlag 是否公众号 8 公众号；24；56
}

// 群成员类型
type RoomMember struct {
	AttrStatus      int
	DisplayName     string
	KeyWord         string
	MemberStatus    string
	NickName        string
	PYInitial       string // 昵称简拼
	PYQuanPin       string // 昵称全拼
	RemarkPYInitial string
	RemarkPYQuanPin string
	Uin             int
	UserName        string
}

// 修改联系人
type ModContact struct {
	Alias             string
	AttrStatus        int
	ChatRoomOwner     string
	City              string
	ContactFlag       int
	ContactType       int
	HeadImgUpdateFlag int // 头像更新标记
	HeadImgUrl        string
	HideInputBarFlag  int
	KeyWord           string
	MemberCount       int
	MemberList        []Member
	NickName          string
	Province          string
	RemarkName        string
	Sex               string
	Signature         string
	SnsFlag           int
	Statues           int
	UserName          string
	VerifyFlag        int
}
