package main

import (
	"expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"github.com/dingdayu/wxbot/handlers/web"
	"github.com/dingdayu/wxbot/handlers/api"
)

const HTTP_IP_PORT = ":8080"

func main() {
	//cron.Test()

	//m := []bson.M{
	//	{"$match": bson.M{"uuid" : "watgL2sJyw=="}},
	//	{"$group": bson.M{"_id": "$province", "count": bson.M{"$sum": 1}}},
	//}
	//provinc := model.GetContactArea(m)
	//fmt.Println(provinc)

	httpGet()
}

func httpGet() {

	http.HandleFunc("/", safeWebHandler(web.Hello))
	http.HandleFunc("/debug/vas", safeWebHandler(metricsHandler))

	http.HandleFunc("/api/send/text", safeWebHandler(api.SendText))
	http.HandleFunc("/api/user/login", safeWebHandler(api.GetUUID))
	http.HandleFunc("/api/user/logout", safeWebHandler(api.Logout))
	http.HandleFunc("/api/user/status", safeWebHandler(api.GetStatus))
	http.HandleFunc("/api/user/all", safeWebHandler(api.GetAllStatus))

	http.HandleFunc("/api/uuid/contact/provinc/list", safeWebHandler(api.GetUuidArea))

	// 监听端口 8080
	fmt.Println("	启动成功！")
	fmt.Println("http://" + HTTP_IP_PORT)
	fmt.Println("------------------------")
	err := http.ListenAndServe(HTTP_IP_PORT, nil)
	if err != nil {
		panic(err)
	}
}

// 服务器内部错误拦截
func safeWebHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 遇到错误时的扫尾工作
		defer func() {
			// 终止（拦截）错误的传递
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50x错误页面
				//w.WriteHeader(http.StatusInternalServerError)
				//handlers.LoadHtml(w, "./templates/50x.html", nil)
				log.Printf("WARN: panic in %v. - %v\n", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		// 调用传入的方法名
		fn(w, r)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(w, "%q: %q", key, str)
		} else {
			fmt.Fprintf(w, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
