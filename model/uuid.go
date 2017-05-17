package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const UUID_COLLECTION_NAME = "uuid"

type UUIDDBT struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	UUID string        `bson:"uuid"`
	Time int           `bson:"time,omitempty"`

	Sex       int    `bson:"sex,omitempty"`
	Signature string `bson:"signature,omitempty"`

	Uin        int    `bson:"uin,omitempty"`
	UserName   string `bson:"username,omitempty"`
	NickName   string `bson:"nickname,omitempty"`
	HeadImgUrl string `bson:"head_img_url,omitempty"`
	Status     string `bson:"status,omitempty"`
}

func UpsertUUID(u UUIDDBT) {
	query := func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"uuid": u.UUID}, bson.M{"$set": u})
		//fmt.Printf("%+v\n", changeInfo)
		return err

	}
	err := witchCollection(UUID_COLLECTION_NAME, query)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetUUIDTByUUID(uuid string) *UUIDDBT {

	uuidT := new(UUIDDBT)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uuid": uuid}).One(&uuidT)
	}
	witchCollection(UUID_COLLECTION_NAME, query)
	return uuidT
}
