package main

import (
	"github.com/dingdayu/wxbot/cron"
)

func main() {
	httpGet()
}

func httpGet() {

	cron.Test()
	//var ts = model.Ts{
	//	UUID: "23",
	//	UserName:"dingdayu",
	//	Time: 123457,
	//}
	//model.UpsertTs(ts)
}
