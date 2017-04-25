package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const USERS_COLLECTION_NAME = "users"

type User struct {
	Id                bson.ObjectId `bson:"_id,omitempty"`
	Uin               int           `bson:"uin,omitempty"`
	UserName          string        `bson:"username,omitempty"`
	NickName          string        `bson:"nickname,omitempty"`
	HeadImgUrl        string        `bson:"head_img_url,omitempty"`
	RemarkName        string        `bson:"remark_name,omitempty"`
	PYInitial         string        `bson:"py_initial,omitempty"`
	PYQuanPin         string        `bson:"py_quan_pin,omitempty"`
	RemarkPYInitial   string        `bson:"remark_py_initial,omitempty"`
	RemarkPYQuanPin   string        `bson:"remark_py_quan_pin,omitempty"`
	HideInputBarFlag  int           `bson:"hide_input_bar_flag,omitempty"`
	StarFriend        int           `bson:"start_friend,omitempty"`
	Sex               int           `bson:"sex,omitempty"`
	Signature         string        `bson:"signature,omitempty"`
	AppAccountFlag    int           `bson:"app_account_flag,omitempty"`
	VerifyFlag        int           `bson:"verify_flag,omitempty"`
	ContactFlag       int           `bson:"contact_flag,omitempty"`
	WebWxPluginSwitch int           `bson:"web_wx_plugin_switch,omitempty"`
	HeadImgFlag       int           `bson:"head_img_flag,omitempty"`
	SnsFlag           int           `bson:"sns_flag,omitempty"`
	UUID              string        `bson:"uuid"`
	Time              int           `bson:"time,omitempty"`
}

/**
 * User
 */
func AddUser(p User) string {
	p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection(USERS_COLLECTION_NAME, query)
	if err != nil {
		return "false"
	}
	return p.Id.Hex()
}

func UpsertUser(u User) {
	query := func(c *mgo.Collection) error {
		c.EnsureIndex(mgo.Index{
			Key:    []string{"uin", "uuid"},
			Unique: true,
			Name:   "uin_uuid",
		})
		changeInfo, err := c.Upsert(bson.M{"uuid": u.UUID, "uin": u.Uin}, bson.M{"$set": u})
		fmt.Printf("%+v\n", changeInfo)
		return err

	}
	err := witchCollection(USERS_COLLECTION_NAME, query)
	if err != nil {
		fmt.Println(err.Error())
	}
}

/**
 * 获取一条记录通过objectid
 */
func GetUserByUIN(uin int) *User {

	user := new(User)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uin": uin}).One(&user)
	}
	witchCollection(USERS_COLLECTION_NAME, query)
	return user
}

//获取所有的person数据
func getAllUser() []User {
	var persons []User
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&persons)
	}
	err := witchCollection(USERS_COLLECTION_NAME, query)
	if err != nil {
		return persons
	}
	return persons
}

//更新person数据
func UpdateUser(query bson.M, change bson.M) bool {
	fmt.Println(change)
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection(USERS_COLLECTION_NAME, exop)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	return false
}
