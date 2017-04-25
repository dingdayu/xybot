package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const CONTACT_COLLECTION_NAME = "contact"

// 联系人
type Contact struct {
	Id               bson.ObjectId `bson:"_id,omitempty"`
	Uin              int           `bson:"uin,omitempty"`
	UserName         string        `bson:"username,omitempty"`
	NickName         string        `bson:"nickname,omitempty"` // 昵称
	HeadImgUrl       string        `bson:"head_img_url,omitempty"`
	ContactFlag      int           `bson:"contact_flag,omitempty"` // 联系人标记，好友1；群2;公众号号3，KeyWord："gh_"为订阅号
	MemberCount      int           `bson:"member_count,omitempty"` // 聊天室成员数量
	MemberList       []RoomMember  `bson:"member_list,omitempty"`  // 聊天室成员
	RemarkName       string        `bson:"remark_name,omitempty"`  // 备注
	HideInputBarFlag int           `bson:"hide_input_bar_flag,omitempty"`
	Sex              int           `bson:"sex,omitempty"`
	Signature        string        `bson:"signature,omitempty"`
	VerifyFlag       int           `bson:"verify_flag,omitempty"` // VerifyFlag 是否公众号 8 公众号；24；56
	OwnerUin         int           `bson:"owner_uin,omitempty"`   // 所有者id
	PYInitial        string        `bson:"py_initial,omitempty"`  // 昵称简拼
	PYQuanPin        string        `bson:"py_quan_pin,omitempty"` // 昵称全拼
	RemarkPYInitial  string        `bson:"remark_py_initial,omitempty"`
	RemarkPYQuanPin  string        `bson:"remark_py_quan_pin,omitempty"`
	StarFriend       int           `bson:"start_friend,omitempty"`
	AppAccountFlag   int           `bson:"app_account_flag,omitempty"`
	Statues          int           `bson:"statues,omitempty"`
	AttrStatus       int           `bson:"attr_status,omitempty"`
	Province         string        `bson:"province,omitempty"` // 省份
	City             string        `bson:"city,omitempty"`     // 城市
	Alias            string        `bson:"alias,omitempty"`    // 别号
	SnsFlag          int           `bson:"sns_flag,omitempty"`
	UniFriend        int           `bson:"uni_friend,omitempty"`
	DisplayName      string        `bson:"display_name,omitempty"`
	ChatRoomId       int           `bson:"chat_romm_id,omitempty"` // 是否聊天室 0 否
	KeyWord          string        `bson:"key_word,omitempty"`
	EncryChatRoomId  string        `bson:"encry_chat_romm_id,omitempty"`
	IsOwner          int           `bson:"is_owner,omitempty"` // 是否所有者

	LoginUin    int    `bson:"login_uin,omitempty"`
	UUID        string `bson:"uuid,omitempty"`
	ContactType string `bson:"contact_type,omitempty"`
}

type RoomMember struct {
	AttrStatus      int    `bson:"attr_status,omitempty"`
	DisplayName     string `bson:"display_name,omitempty"`
	KeyWord         string `bson:"key_word,omitempty"`
	MemberStatus    int    `bson:"member_status,omitempty"`
	NickName        string `bson:"nickname,omitempty"`
	PYInitial       string `bson:"py_initial,omitempty"`  // 昵称简拼
	PYQuanPin       string `bson:"py_quan_pin,omitempty"` // 昵称全拼
	RemarkPYInitial string `bson:"remark_py_initial,omitempty"`
	RemarkPYQuanPin string `bson:"remark_py_quan_pin,omitempty"`
	Uin             int    `bson:"uin,omitempty"`
	UserName        string `bson:"username,omitempty"`
}

/**
 * User
 */
func AddContact(p Contact) string {
	p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection(CONTACT_COLLECTION_NAME, query)
	if err != nil {
		return "false"
	}
	return p.Id.Hex()
}

func UpsertContact(p *Contact) {
	query := func(c *mgo.Collection) error {
		changeInfo, err := c.Upsert(bson.M{"uuid": p.UUID, "username": p.UserName}, bson.M{"$set": p})
		fmt.Printf("%+v\n", changeInfo)
		return err

	}
	err := witchCollection(CONTACT_COLLECTION_NAME, query)
	if err != nil {
		fmt.Println(err.Error())
	}
}

/**
 * 获取一条记录通过objectid
 */
func GetContactById(id string) *Contact {
	objid := bson.ObjectIdHex(id)
	person := new(Contact)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&person)
	}
	witchCollection(CONTACT_COLLECTION_NAME, query)
	return person
}

//获取所有的person数据
func PageContact() []Contact {
	var persons []Contact
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&persons)
	}
	err := witchCollection(CONTACT_COLLECTION_NAME, query)
	if err != nil {
		return persons
	}
	return persons
}

func GetChatRoomContact() []Contact {
	var contact []Contact
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&contact)
	}
	err := witchCollection(CONTACT_COLLECTION_NAME, query)
	if err != nil {
		return contact
	}
	return contact
}

//更新person数据
func UpdateContact(query bson.M, change bson.M) bool {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection(CONTACT_COLLECTION_NAME, exop)
	if err != nil {
		return true
	}
	return false
}
