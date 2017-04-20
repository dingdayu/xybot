package types

type User struct {
	HeadImgUrl        string
	UserName          string
	Uit               int
	NickName          string
	RemarkNam         string
	Sex               int
	Signature         string
	SnsFlag           int
	StarFriend        int
	VerifyFlag        int
	WebWxPluginSwitch int
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

var SPECIAL_USERS = []string{"newsapp", "fmessage", "filehelper", "weibo", "qqmail",
	"fmessage", "tmessage", "qmessage", "qqsync", "floatbottle",
	"lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp",
	"blogapp", "facebookapp", "masssendapp", "meishiapp",
	"feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder",
	"weixinreminder", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c",
	"officialaccounts", "notification_messages", "wxid_novlwrv3lqwv11",
	"gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages"}
