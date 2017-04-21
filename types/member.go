package types

// 好友信息类型
type Member struct {
	Uin              int
	UserName         string
	NickName         string // 昵称
	HeadImgUrl       string
	ContactFlag      int        // 联系人标记，好友1；群2
	MemberCount      int        // 聊天室成员数量
	MemberList       RoomMember // 聊天室成员
	RemarkName       string     // 备注
	HideInputBarFlag int
	Sex              int
	Signature        string
	VerifyFlag       int    // VerifyFlag 是否公众号 8 公众号；24；56
	OwnerUin         int    // 所有者id
	PYInitial        string // 昵称简拼
	PYQuanPin        string // 昵称全拼
	RemarkPYInitial  string
	RemarkPYQuanPin  string
	StarFriend       int
	AppAccountFlag   int
	Statues          int
	AttrStatus       int
	Province         string // 省份
	City             string // 城市
	Alias            string // 别号
	SnsFlag          int
	UniFriend        int
	DisplayName      string
	ChatRoomId       int // 是否聊天室 0 否
	KeyWord          string
	EncryChatRoomId  string
	IsOwner          int // 是否所有者
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
