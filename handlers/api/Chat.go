package api

import (
	"github.com/dingdayu/wxbot/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// 获取省和城市的计数
func GetUuidArea(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	ret := RetT{}
	if uuid == "" {
		ret = RetT{Code: 302, Msg: "uuid Not Empty"}
		RetJson(ret, w)
		return
	}
	m := []bson.M{
		{"$match": bson.M{"uuid": uuid, "contact_type": "Friends"}},
		{"$group": bson.M{"_id": "$province", "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"name": "$_id", "value": "$count", "_id": false}},
	}
	provinc := model.GetContactArea(m)

	m = []bson.M{
		{"$match": bson.M{"uuid": uuid, "contact_type": "Friends"}},
		{"$group": bson.M{"_id": "$city", "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"name": "$_id", "value": "$count", "_id": false}},
	}
	city := model.GetContactArea(m)

	data := make(map[string]interface{})
	data["provinc"] = provinc
	data["city"] = city
	ret = RetT{Code: 200, Msg: "success", Data: data}
	RetJson(ret, w)
}
