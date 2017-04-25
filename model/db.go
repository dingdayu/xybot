package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const URL = "" //mongodb连接字符串

var (
	mgoSession *mgo.Session
	dataBase   = "mydb"
)

func init() {
	//initExindex()
}

//user
// 联系人
// 群
// 公众号
// 公众号文章

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

/**
 * 执行查询，此方法可拆分做为公共方法
 * [SearchPerson description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           bson.M [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func SearchPerson(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = witchCollection(collectionName, exop)
	return
}

func initExindex() {
	//getSession().DB(dataBase).C(USERS_COLLECTION_NAME).EnsureIndex(mgo.Index{
	//	Key:    []string{"uin", "uuid"},
	//	Unique: true,
	//	Name:   "uin_uuid",
	//})
	//getSession().DB(dataBase).C(CONTACT_COLLECTION_NAME).EnsureIndex(mgo.Index{
	//	Key:    []string{"login_uin", "uuid"},
	//	Unique: true,
	//	Name:   "uin_uuid",
	//})
}
