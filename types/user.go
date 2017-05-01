package types

type User struct {
	Uin               int
	UserName          string
	NickName          string
	HeadImgUrl        string
	RemarkName        string
	PYInitial         string
	PYQuanPin         string
	RemarkPYInitial   string
	RemarkPYQuanPin   string
	HideInputBarFlag  int
	StarFriend        int
	Sex               int
	Signature         string
	AppAccountFlag    int
	VerifyFlag        int
	ContactFlag       int
	WebWxPluginSwitch int
	HeadImgFlag       int
	SnsFlag           int
}

type Profile struct {
	Alias     string
	BindEmail struct {
		Buff string
	}
	BindMobile struct {
		Buff string
	}
	BindUin           int
	BitFlag           int
	HeadImgUpdateFlag int
	HeadImgUrl        string
	NickName          struct {
		Buff string
	}
	PersonalCard int
	Sex          int
	Signature    string
	Status       int
	UserName     struct {
		Buff string
	}
}

var SPECIAL_USERS = []string{"newsapp", "filehelper", "weibo", "qqmail",
	"fmessage", "tmessage", "qmessage", "qqsync", "floatbottle",
	"lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp",
	"blogapp", "facebookapp", "masssendapp", "meishiapp",
	"feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder",
	"weixinreminder", "wxid_novlwrv3lqwv11",
	"officialaccounts",
	"gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages", "notifymessage"}

var SPECIAL_USERS_NAME = []string{
	"微信运动",
}

var ENTERPRISE_NAME = []string{
	"微信摇一摇周边",
}
